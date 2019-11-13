package simplelb

import (
	"log"
	"net/url"
	"strings"
)

// CreateServerPool creates a server pool by serverList
func CreateServerPool(serverList string) *ServerPool {
	var serverPool ServerPool

	index := 0

	for _, tok := range strings.Split(serverList, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		serverURL, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		b := &Backend{URL: serverURL, Alive: true}
		serverPool.addBackend(b)
		b.Proxy = serverPool.createProxy(index, serverURL)
		index++
		log.Printf("Configured server: %s\n", serverURL)
	}

	return &serverPool
}
