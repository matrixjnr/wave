package main

import (
	"github.com/matrixjnr/wave/internal/channels"
	"github.com/matrixjnr/wave/internal/routing"
	"github.com/matrixjnr/wave/pkg/quic"
	"log"
)

func main() {
	// Initialize the router
	router := routing.NewRouter()

	// Start QUIC server
	server := quic.NewServer("0.0.0.0:4242")

	// Start listening for sessions
	go func() {
		if err := server.ListenAndServe(router); err != nil {
			log.Fatalf("QUIC Server error: %v", err)
		}
	}()

	log.Println("Wave Protocol QUIC server started on 0.0.0.0:4242")

	// Simulated subscription for internal messages
	internalListener := router.Subscribe("server-channel")
	go func() {
		for msg := range internalListener {
			log.Printf("Internal message received: %+v", msg)
		}
	}()

	// Keep the server alive
	select {}
}
