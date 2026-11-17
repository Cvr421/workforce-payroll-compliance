package main

import (
	"context"
	"fmt"
	"os"

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

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		log.WithError(err).Error("failed to connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("seed command started - implementation in Phase 9")
	_ = pool
	_ = cfg.SeedPassword
}