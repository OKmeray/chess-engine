package env_test

import (
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/env"
)

func TestGetString(t *testing.T) {
	t.Run("returns fallback when unset", func(t *testing.T) {
		const (
			key      = "UNSET_VAR"
			fallback = "default_val"
		)
		got := env.GetString(key, fallback)
		if got != fallback {
			t.Errorf("GetString(%q) = %q, want %q", key, got, fallback)
		}
	})

	t.Run("returns env value when set", func(t *testing.T) {
		const (
			key      = "TEST_VAR"
			val      = "custom_val"
			fallback = "default_val"
		)
		t.Setenv(key, val)
		got := env.GetString(key, fallback)
		if got != val {
			t.Errorf("GetString(%q) = %q, want %q", key, got, val)
		}
	})
}

func TestGetInt(t *testing.T) {
	t.Run("returns fallback when unset", func(t *testing.T) {
		const (
			key      = "UNSET_INT"
			fallback = 21
		)
		val, err := env.GetInt(key, fallback)
		if err != nil {
			t.Fatalf("GetInt(%q): unexpected error: %v", key, err)
		}
		if val != fallback {
			t.Errorf("GetInt(%q) = %d, want %d", key, val, fallback)
		}
	})

	t.Run("returns parsed int", func(t *testing.T) {
		const (
			key      = "TEST_INT"
			fallback = 21
			want     = 100
		)
		t.Setenv(key, "100")
		val, err := env.GetInt(key, fallback)
		if err != nil {
			t.Fatalf("GetInt(%q): unexpected error: %v", key, err)
		}
		if val != want {
			t.Errorf("GetInt(%q) = %d, want %d", key, val, want)
		}
	})

	t.Run("returns error on invalid int", func(t *testing.T) {
		const (
			key      = "TEST_INT"
			fallback = 21
		)
		t.Setenv(key, "not_a_number")
		_, err := env.GetInt(key, fallback)
		if err == nil {
			t.Errorf("GetInt(%q) error = %v, wantErr %v", key, err, true)
		}
	})
}

func TestGetDuration(t *testing.T) {
	t.Run("returns fallback when unset", func(t *testing.T) {
		const (
			key      = "UNSET_DUR"
			fallback = 5 * time.Second
		)
		val, err := env.GetDuration(key, fallback)
		if err != nil {
			t.Fatalf("GetDuration(%q): unexpected error: %v", key, err)
		}
		if val != fallback {
			t.Errorf("GetDuration(%q) = %v, want %v", key, val, fallback)
		}
	})

	t.Run("returns parsed duration", func(t *testing.T) {
		const (
			key      = "TEST_DUR"
			fallback = time.Second
			want     = 10 * time.Second
		)
		t.Setenv(key, "10s")
		val, err := env.GetDuration(key, fallback)
		if err != nil {
			t.Fatalf("GetDuration(%q): unexpected error: %v", key, err)
		}
		if val != want {
			t.Errorf("GetDuration(%q) = %v, want %v", key, val, want)
		}
	})

	t.Run("returns error on invalid duration", func(t *testing.T) {
		const (
			key      = "TEST_DUR"
			fallback = time.Second
		)
		t.Setenv(key, "invalid_duration")
		_, err := env.GetDuration(key, fallback)
		if err == nil {
			t.Errorf("GetDuration(%q) error = %v, wantErr %v", key, err, true)
		}
	})
}
