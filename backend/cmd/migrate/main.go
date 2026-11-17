package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/complyhub/backend/internal/config"
	"github.com/complyhub/backend/internal/db"
	"github.com/complyhub/backend/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("dir", "up", "migration direction: up or down")
	steps := flag.Int("steps", 0, "number of steps to migrate (0 = all)")
	flag.Parse()

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

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.WithError(err).Error("failed to create migrate instance")
		os.Exit(1)
	}

	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid direction: %s\n", *direction)
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.WithError(err).Error("migration failed")
		os.Exit(1)
	}

	log.Info("migration completed", "direction", *direction)
}