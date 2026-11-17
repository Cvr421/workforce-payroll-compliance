package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	JWTIssuer       string
	JWTAudience     string
	Port            string
	LogLevel        string
	CORSAllowedOrigins string
	StorageDir      string
	DBMaxConns      int
	WorkerConcurrency int
	WorkerPollInterval time.Duration
	JobStaleAfter   time.Duration
	SeedPassword    string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisURL:           getEnv("REDIS_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTIssuer:          getEnv("JWT_ISSUER", "complylite"),
		JWTAudience:        getEnv("JWT_AUDIENCE", "complylite-api"),
		Port:               getEnv("PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		StorageDir:         getEnv("STORAGE_DIR", "/tmp/complylite-storage"),
		DBMaxConns:         getEnvInt("DB_MAX_CONNS", 10),
		WorkerConcurrency:  getEnvInt("WORKER_CONCURRENCY", 4),
		WorkerPollInterval: getEnvDuration("WORKER_POLL_INTERVAL", 2*time.Second),
		JobStaleAfter:      getEnvDuration("JOB_STALE_AFTER", 5*time.Minute),
		SeedPassword:       getEnv("SEED_PASSWORD", "DemoPassword123"),
	}

	if cfg.DatabaseURL == "" {
		return nil, ErrMissingEnv("DATABASE_URL")
	}
	if cfg.JWTSecret == "" {
		return nil, ErrMissingEnv("JWT_SECRET")
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, ErrInvalidEnv("JWT_SECRET", "must be at least 32 bytes")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultValue
}

type ConfigError struct {
	Key    string
	Reason string
}

func (e *ConfigError) Error() string {
	return "config: " + e.Key + ": " + e.Reason
}

func ErrMissingEnv(key string) error {
	return &ConfigError{Key: key, Reason: "required but not set"}
}

func ErrInvalidEnv(key, reason string) error {
	return &ConfigError{Key: key, Reason: reason}
}