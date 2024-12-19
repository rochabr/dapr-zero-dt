package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	daprClient dapr.Client
	appPort    = "8080"
	healthPort = "3000" // Separate port for health checks
)

func main() {
	var err error
	daprClient, err = dapr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer daprClient.Close()

	// Start the ping loop in a goroutine
	go pingLoop()

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
			if daprClient != nil {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "Ready")
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
		})

		log.Printf("Health server listening on port %s", healthPort)
		if err := http.ListenAndServe(":"+healthPort, mux); err != nil {
			log.Fatalf("Health server failed: %v", err)
		}
	}()

	// Start main application server
	appMux := http.NewServeMux()
	log.Printf("Ping service listening on port %s", appPort)
	if err := http.ListenAndServe(":"+appPort, appMux); err != nil {
		log.Fatal(err)
	}
}

func pingLoop() {
	for {
		ctx := context.Background()
		content := &dapr.DataContent{
			ContentType: "text/plain",
			Data:        []byte("ping"),
		}

		resp, err := daprClient.InvokeMethodWithContent(ctx, "pong-service", "ping", "post", content)
		if err != nil {
			log.Printf("Error invoking pong service: %v", err)
		} else {
			log.Printf("Received response: %s", string(resp))
		}

		time.Sleep(1 * time.Second)
	}
}
