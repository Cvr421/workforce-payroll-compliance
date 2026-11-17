package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/complyhub/backend/internal/config"
	"github.com/complyhub/backend/internal/db"
	"github.com/complyhub/backend/internal/httpx"
	"github.com/complyhub/backend/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
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

	server := httpx.NewServer(cfg, log)

	// Setup routes will be added in Phase 2
	server.Router().HandleFunc("GET /healthz", server.Healthz)
	server.Router().HandleFunc("GET /readyz", server.Readyz)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("server failed")
			os.Exit(1)
		}
	}()

	log.Info("server started", "port", cfg.Port)

	<-ctx.Done()
	log.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("server shutdown failed")
		os.Exit(1)
	}

	log.Info("server stopped")
}