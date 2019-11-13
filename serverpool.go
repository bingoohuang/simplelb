package simplelb

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
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
	if s.backendsNum == 1 {
		return s.backends[0]
	}

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
		alive := IsAddressAlive(b.Host)

		if oldAlive == alive {
			continue
		}

		b.SetAlive(alive)

		status := "up"
		if !alive {
			status = "down"
		}

		log.Printf("%s [%s]\n", b.Host, status)
	}
}

// Lb load balances the incoming request
func (s *ServerPool) Lb(ctx *fasthttp.RequestCtx) {
	if peer := s.GetNextPeer(); peer != nil {
		peer.Proxy.ServeHTTP(ctx)
		return
	}

	ctx.Error("Service not available", http.StatusServiceUnavailable)
}

// HealthCheck runs a routine for check status of the backends every 20s
func (s *ServerPool) HealthCheck() {
	for range time.NewTicker(time.Second * 20).C {
		s.healthCheck()
	}
}
