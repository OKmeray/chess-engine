package config

import (
	"fmt"
	"time"

	"github.com/OKmeray/chess-engine/internal/env"
)

// Config holds the complete configuration for the application
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Engine EngineConfig
}

// ServerConfig defines the HTTP server network and timeout parameters.
type ServerConfig struct {
	Addr            string
	Env             string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	CacheTTL        time.Duration
}

// DBConfig defines connection pool parameters for PostgreSQL.
type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// EngineConfig defines resource limits for the search engine.
type EngineConfig struct {
	MaxNodes   int
	NumWorkers int
}

// Load reads environment variables and constructs the complete application configuration.
func Load() (Config, error) {
	dbUser := env.GetString("POSTGRES_USER", "postgres")
	dbPass := env.GetString("POSTGRES_PASSWORD", "postgres")
	dbHost := env.GetString("POSTGRES_HOST", "localhost")
	dbPort := env.GetString("POSTGRES_PORT", "5433")
	dbName := env.GetString("POSTGRES_DB", "chess_engine")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	maxOpenConns, err := env.GetInt("DB_MAX_OPEN_CONNS", 30)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}

	maxIdleConns, err := env.GetInt("DB_MAX_IDLE_CONNS", 30)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	maxIdleTime, err := env.GetDuration("DB_MAX_IDLE_TIME", time.Minute*15)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_MAX_IDLE_TIME: %w", err)
	}

	addr := ":" + env.GetString("PORT", "8080")
	srvEnv := env.GetString("ENV", "development")

	writeTimeout, err := env.GetDuration("WRITE_TIMEOUT", time.Second*30)
	if err != nil {
		return Config{}, fmt.Errorf("invalid WRITE_TIMEOUT: %w", err)
	}

	readTimeout, err := env.GetDuration("READ_TIMEOUT", time.Second*10)
	if err != nil {
		return Config{}, fmt.Errorf("invalid READ_TIMEOUT: %w", err)
	}

	idleTimeout, err := env.GetDuration("IDLE_TIMEOUT", time.Minute*1)
	if err != nil {
		return Config{}, fmt.Errorf("invalid IDLE_TIMEOUT: %w", err)
	}

	shutdownTimeout, err := env.GetDuration("SHUTDOWN_TIMEOUT", time.Minute*1)
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	cacheTTL, err := env.GetDuration("CACHE_TTL", time.Minute*5)
	if err != nil {
		return Config{}, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}

	maxNodes, err := env.GetInt("MCTS_MAX_NODES", 1000000)
	if err != nil {
		return Config{}, fmt.Errorf("invalid MCTS_MAX_NODES: %w", err)
	}

	numWorkers, err := env.GetInt("MCTS_NUM_WORKERS", 8)
	if err != nil {
		return Config{}, fmt.Errorf("invalid MCTS_NUM_WORKERS: %w", err)
	}

	return Config{
		Server: ServerConfig{
			Addr:            addr,
			Env:             srvEnv,
			WriteTimeout:    writeTimeout,
			ReadTimeout:     readTimeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
			CacheTTL:        cacheTTL,
		},
		DB: DBConfig{
			DSN:          dsn,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxIdleTime:  maxIdleTime,
		},
		Engine: EngineConfig{
			MaxNodes:   maxNodes,
			NumWorkers: numWorkers,
		},
	}, nil
}
