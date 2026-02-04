package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pc-stats-api/internal/api"
	"pc-stats-api/internal/collector"
	"pc-stats-api/internal/config"
	"pc-stats-api/internal/storage"
)

func main() {
	log.Println("Starting Worker Monitoring Agent...")

	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration: Port=%s, Interval=%ds, HistorySize=%d",
		cfg.Port, cfg.IntervalSec, cfg.HistorySize)

	// Initialize components
	systemCollector := collector.NewSystemCollector()
	ringBuffer := storage.NewRingBuffer(cfg.HistorySize)

	// Start metrics collection goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go collectMetrics(ctx, systemCollector, ringBuffer, cfg.IntervalSec)

	// Start HTTP server (with embedded web files)
	server := api.NewServer(ringBuffer, cfg.IntervalSec, GetWebFS())

	// Setup graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		cancel()
		os.Exit(0)
	}()

	// Start server (blocking)
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Web UI available at http://localhost:%s/ui/", cfg.Port)
	if err := server.Start(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func collectMetrics(ctx context.Context, collector *collector.SystemCollector, storage *storage.RingBuffer, intervalSec int) {
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	// Collect initial sample immediately
	if sample, err := collector.Collect(); err == nil {
		storage.Add(sample)
		log.Printf("Initial metrics collected: CPU=%.1f%%, RAM=%.1f%%",
			sample.CPU.Usage*100, sample.RAM.Usage*100)
	} else {
		log.Printf("Failed to collect initial metrics: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Metrics collection stopped")
			return
		case <-ticker.C:
			sample, err := collector.Collect()
			if err != nil {
				log.Printf("Error collecting metrics: %v", err)
				continue
			}

			storage.Add(sample)

			// Log summary
			gpuInfo := ""
			if sample.GPU != nil {
				gpuInfo = fmt.Sprintf(", GPU=%.1f%%", sample.GPU.Util*100)
			}
			log.Printf("Metrics: CPU=%.1f%%, RAM=%.1f%%%s",
				sample.CPU.Usage*100, sample.RAM.Usage*100, gpuInfo)
		}
	}
}
