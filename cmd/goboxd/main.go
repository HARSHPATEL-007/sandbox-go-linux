package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"goboxd/internal/api"
	"goboxd/internal/config"
	"goboxd/internal/engine"
)

func main() {
	log.Println("Starting goboxd...")

	// 1. Load Registry
	registry, err := config.Load("config/languages.yaml")
	if err != nil {
		log.Fatalf("Failed to load languages: %v", err)
	}
	log.Printf("Loaded %d languages from registry", len(registry))

	// 2. Start Engine with Bounded Concurrency
	concurrency := runtime.NumCPU()
	if envC := os.Getenv("MAX_CONCURRENCY"); envC != "" {
		fmt.Sscanf(envC, "%d", &concurrency)
	}
	eng := engine.NewEngine(concurrency)
	log.Printf("Engine started with %d concurrent workers", concurrency)

	// 3. Start API
	router := api.NewRouter(registry, eng)
	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}