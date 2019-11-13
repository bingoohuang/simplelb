package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bingoohuang/simplelb"
)

func main() {
	backends := ""
	port := 0

	flag.StringVar(&backends, "b", "", "Load balanced backends, use , to separate")
	flag.IntVar(&port, "p", 3030, "Port to serve")
	flag.Parse()

	if backends == "" {
		log.Fatal("Please provide one or more backends to load balance")
	}

	serverPool := simplelb.CreateServerPool(backends)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(serverPool.Lb),
	}

	go serverPool.HealthCheck()

	log.Printf("Load Balancer started at :%d\n", port)

	log.Fatal(server.ListenAndServe())
}
