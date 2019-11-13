package simplelb

import (
	"sync"
)

// Backend holds the data about a server
type Backend struct {
	Alive bool
	IsTLS bool
	Host  string // ip:port
	mux   sync.RWMutex
	Proxy *ReverseProxy
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()

	return
}
