package simplelb

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

// ServerPool holds information about reachable backends
type ServerPool struct {
	backends    []*Backend
	backendsNum int
	current     uint64
}

// addBackend to the server pool
func (s *ServerPool) addBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
	s.backendsNum++
}

// nextIndex atomically increase the counter and return an index
func (s *ServerPool) nextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(s.backendsNum))
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.nextIndex()     // loop entire backends to find out an Alive backend
	l := s.backendsNum + next // start from next and move a full cycle

	for i := next; i < l; i++ {
		idx := i % s.backendsNum       // take an index by modding
		if s.backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}

			return s.backends[idx]
		}
	}

	return nil
}

// healthCheck pings the backends and update the status
func (s *ServerPool) healthCheck() {
	for _, b := range s.backends {
		oldAlive := b.IsAlive()
		alive := IsAddressAlive(b.URL.Host)

		if oldAlive == alive {
			continue
		}

		b.SetAlive(alive)

		status := "up"
		if !alive {
			status = "down"
		}

		log.Printf("%s [%s]\n", b.URL, status)
	}
}

// Lb load balances the incoming request
func (s *ServerPool) Lb(w http.ResponseWriter, r *http.Request) {
	if attempts := GetAttempts(r); attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)

		return
	}

	if peer := s.GetNextPeer(); peer != nil {
		peer.Proxy.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// HealthCheck runs a routine for check status of the backends every 20s
func (s *ServerPool) HealthCheck() {
	for range time.NewTicker(time.Second * 20).C {
		s.healthCheck()
	}
}

func (s *ServerPool) createProxy(backIndex int, backURL *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(backURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		s.retry(backIndex, proxy, w, r, e)
	}
	proxy.Transport = &http.Transport{
		DialContext:           TimeoutDialContext(60*time.Second, 60*time.Second),
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return proxy
}

func (s *ServerPool) retry(backIndex int, p http.Handler, w http.ResponseWriter, r *http.Request, e error) {
	backend := s.backends[backIndex]
	log.Printf("ErrorHandler [%s] %s\n", backend.URL.Host, e.Error())

	if retries := GetRetry(r); retries < 3 {
		<-time.After(10 * time.Millisecond)

		ctx := context.WithValue(r.Context(), Retries, retries+1)
		p.ServeHTTP(w, r.WithContext(ctx))

		return
	}

	// after 3 retries, mark this backend as down
	backend.SetAlive(false)

	// if the same r routing for few attempts with different backends, increase the count
	attempts := GetAttempts(r)
	log.Printf("%s(%s) Attempting retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
	ctx := context.WithValue(r.Context(), Attempts, attempts+1)
	s.Lb(w, r.WithContext(ctx))
}
