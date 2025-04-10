package main

import (
	"log"
	"net/http"
	"time"

	"nats_prometheus_exporter/config"
	"nats_prometheus_exporter/internal/server"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	collector := &server.Collector{
		Cache:     make(map[string]server.CacheEntry),
		TTL:       300, // evict metrics if not updated in 5 min
		CleanupHz: 1 * time.Minute,
	}

	err = collector.Start(cfg.NatsURL, cfg.Subject)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}

	log.Printf("Serving metrics on port %s...", cfg.ListenPort)
	http.Handle("/metrics", collector)
	if err := http.ListenAndServe(":"+cfg.ListenPort, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
