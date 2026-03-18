package main

import (
	"log"

	"github.com/joekariuki/sandgrouse/internal/proxy"
)

func main() {
	srv := &proxy.Server{
		ListenAddr: ":8080",
	}
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start proxy: %v", err)
	}
}
