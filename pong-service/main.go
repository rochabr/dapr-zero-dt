package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	appPort    = "8080"
	healthPort = "3000" // Separate port for health checks
)

func main() {
	// Start health check server in a separate goroutine
	go func() {
		mux := http.NewServeMux()

		// Health check endpoint
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Healthy")
		})

		// Readiness check endpoint
		mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Ready")
		})

		log.Printf("Health server listening on port %s", healthPort)
		if err := http.ListenAndServe(":"+healthPort, mux); err != nil {
			log.Fatalf("Health server failed: %v", err)
		}
	}()

	// Main application server
	appMux := http.NewServeMux()

	// Ping endpoint
	appMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received ping request")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong from %s", os.Getenv("POD_NAME"))
	})

	log.Printf("Pong service listening on port %s", appPort)
	if err := http.ListenAndServe(":"+appPort, appMux); err != nil {
		log.Fatal(err)
	}
}
