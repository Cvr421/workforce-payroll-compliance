package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/complyhub/backend/internal/config"
	"github.com/complyhub/backend/internal/db"
	"github.com/complyhub/backend/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		log.WithError(err).Error("failed to connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("worker started",
		"concurrency", cfg.WorkerConcurrency,
		"poll_interval", cfg.WorkerPollInterval,
		"job_stale_after", cfg.JobStaleAfter,
	)

	// Worker implementation will be added in Phase 8
	ticker := time.NewTicker(cfg.WorkerPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("worker shutting down...")
			return
		case <-ticker.C:
			_ = pool
		}
	}
}