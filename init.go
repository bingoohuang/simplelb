package simplelb

import (
	"log"
	"strings"
)

// CreateServerPool creates a server pool by serverList
func CreateServerPool(serverList string) *BackendPool {
	var serverPool BackendPool

	for _, tok := range strings.Split(serverList, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		serverURL := tok

		isTLS, host, err := ParseAddress(serverURL)
		if err != nil {
			log.Fatal(err)
		}

		b := &Backend{IsTLS: isTLS, Host: host, Alive: true, Proxy: NewReverseProxy(isTLS, host)}
		serverPool.Add(b)

		log.Printf("Configured server: %s\n", serverURL)
	}

	return &serverPool
}
