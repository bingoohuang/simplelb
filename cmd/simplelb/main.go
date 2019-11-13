package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bingoohuang/simplelb"

	"runtime/pprof"

	"github.com/valyala/fasthttp"
)

func main() {
	cpuProfile()

	backends := ""
	port := 0

	flag.StringVar(&backends, "b", "", "Load balanced backends, use , to separate")
	flag.IntVar(&port, "p", 3030, "Port to serve")
	flag.Parse()

	if backends == "" {
		log.Fatal("Please provide one or more backends to load balance")
	}

	serverPool := simplelb.CreateServerPool(backends)
	addr := fmt.Sprintf(":%d", port)

	go serverPool.HealthCheck()

	log.Printf("Load Balancer started at :%d\n", port)

	log.Fatal(fasthttp.ListenAndServe(addr, serverPool.Lb))
}

func cpuProfile() {
	cpuProfile, _ := os.Create("cpu_profile")
	_ = pprof.StartCPUProfile(cpuProfile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Block until a signal is received.
		s := <-c
		fmt.Println("Got signal:", s)
		pprof.StopCPUProfile()
		_ = cpuProfile.Close()
		os.Exit(-1)
	}()
}
