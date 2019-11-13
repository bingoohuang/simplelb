package simplelb

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

// BackendPool holds information about reachable backends
type BackendPool struct {
	backends []*Backend
	total    uint64
	current  uint64
}

// CheckBackends check backends
func (s *BackendPool) CheckBackends() {
	if s.total == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}
}

// Add adds a backend to the server pool
func (s *BackendPool) Add(backend *Backend) {
	s.backends = append(s.backends, backend)
	s.total++
}

// nextIndex atomically increase the counter and return an index
func (s *BackendPool) nextIndex() uint64 {
	return atomic.AddUint64(&s.current, 1) % s.total
}

// GetNextPeer returns next active peer to take a connection
func (s *BackendPool) GetNextPeer() *Backend {
	if s.total == 1 {
		return s.backends[0]
	}

	next := s.nextIndex() // loop entire backends to find out an Alive backend

	for i := next; i < next+s.total; i++ {
		idx := i % s.total

		if !s.backends[idx].IsAlive() {
			continue
		}

		// if we have an alive backend, use it and store if its not the original one
		if i != next {
			atomic.StoreUint64(&s.current, idx)
		}

		return s.backends[idx]
	}

	// 如果全部下线，则选择第一个进行尝试
	return s.backends[next]
}

// healthCheck pings the backends and update the status
func (s *BackendPool) healthCheck() {
	for _, b := range s.backends {
		oldAlive := b.IsAlive()
		alive := IsAddressAlive(b.Host)

		if oldAlive == alive {
			continue
		}

		b.SetAlive(alive)

		log.Printf("%s alive %v\n", b.Host, alive)
	}
}

// Lb load balances the incoming request
func (s *BackendPool) Lb(ctx *fasthttp.RequestCtx) {
	if peer := s.GetNextPeer(); peer != nil {
		if err := peer.Proxy.ServeHTTP(ctx); err == nil {
			return
		}
	}

	ctx.Error("Service not available", http.StatusServiceUnavailable)
}

// HealthCheck runs a routine for check status of the backends every 20s
func (s *BackendPool) HealthCheck() {
	for range time.NewTicker(time.Second * 20).C {
		s.healthCheck()
	}
}
