package config_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load(): unexpected error: %v", err)
	}

	const (
		wantAddr    = ":8080"
		wantMaxOpen = 30
		wantWorkers = 8
		wantDSN     = "postgres://postgres:postgres@localhost:5433/chess_engine?sslmode=disable"
	)

	if cfg.Server.Addr != wantAddr {
		t.Errorf("Load() Addr = %v, want %v", cfg.Server.Addr, wantAddr)
	}
	if cfg.DB.MaxOpenConns != wantMaxOpen {
		t.Errorf("Load() MaxOpenConns = %v, want %v", cfg.DB.MaxOpenConns, wantMaxOpen)
	}
	if cfg.Engine.NumWorkers != wantWorkers {
		t.Errorf("Load() NumWorkers = %v, want %v", cfg.Engine.NumWorkers, wantWorkers)
	}
	if cfg.DB.DSN != wantDSN {
		t.Errorf("Load() DSN = %v, want %v", cfg.DB.DSN, wantDSN)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	const (
		port        = "9090"
		workers     = 16
		readTimeout = 15 * time.Second
		user        = "custom_user"
		db          = "custom_db"
		wantAddr    = ":" + port
	)
	wantDSN := fmt.Sprintf("postgres://%s:postgres@localhost:5433/%s?sslmode=disable", user, db)

	t.Setenv("PORT", port)
	t.Setenv("MCTS_NUM_WORKERS", strconv.Itoa(workers))
	t.Setenv("READ_TIMEOUT", readTimeout.String())
	t.Setenv("POSTGRES_USER", user)
	t.Setenv("POSTGRES_DB", db)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load(): unexpected error: %v", err)
	}

	if cfg.Server.Addr != wantAddr {
		t.Errorf("Load() Addr = %v, want %v", cfg.Server.Addr, wantAddr)
	}
	if cfg.Engine.NumWorkers != workers {
		t.Errorf("Load() NumWorkers = %v, want %v", cfg.Engine.NumWorkers, workers)
	}
	if cfg.Server.ReadTimeout != readTimeout {
		t.Errorf("Load() ReadTimeout = %v, want %v", cfg.Server.ReadTimeout, readTimeout)
	}
	if cfg.DB.DSN != wantDSN {
		t.Errorf("Load() DSN = %v, want %v", cfg.DB.DSN, wantDSN)
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{key: "DB_MAX_OPEN_CONNS", value: "invalid_int"},
		{key: "DB_MAX_IDLE_CONNS", value: "invalid_int"},
		{key: "DB_MAX_IDLE_TIME", value: "invalid_duration"},
		{key: "WRITE_TIMEOUT", value: "invalid_duration"},
		{key: "READ_TIMEOUT", value: "invalid_duration"},
		{key: "IDLE_TIMEOUT", value: "invalid_duration"},
		{key: "SHUTDOWN_TIMEOUT", value: "invalid_duration"},
		{key: "CACHE_TTL", value: "invalid_duration"},
		{key: "MCTS_MAX_NODES", value: "invalid_int"},
		{key: "MCTS_NUM_WORKERS", value: "invalid_int"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)

			_, err := config.Load()
			if err == nil {
				t.Errorf("Load(%s=%q) error = %v, wantErr %v", tt.key, tt.value, err, true)
			}
		})
	}
}
