package cache_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain"
	"github.com/OKmeray/chess-engine/internal/infrastructure/cache"
)

type mockProvider struct {
	getByIDCalls int
	models       map[int]domain.NNModel
}

func (m *mockProvider) GetByID(ctx context.Context, id int) (domain.NNModel, error) {
	m.getByIDCalls++
	model, ok := m.models[id]
	if !ok {
		return domain.NNModel{}, domain.ErrModelNotFound
	}
	return model, nil
}

func TestModelCache_GetByID(t *testing.T) {
	ctx := context.Background()
	mock := &mockProvider{
		models: map[int]domain.NNModel{
			1: {ID: 1, Name: "ModelA"},
		},
	}

	c := cache.NewModelCache(mock, 1*time.Minute)

	got1, err := c.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID(1): unexpected error: %v", err)
	}
	if got1.Name != "ModelA" {
		t.Errorf("GetByID(1).Name = %q, want %q", got1.Name, "ModelA")
	}
	if mock.getByIDCalls != 1 {
		t.Errorf("GetByID() calls = %d, want %d", mock.getByIDCalls, 1)
	}

	got2, err := c.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID(1): unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got1, got2) {
		t.Errorf("GetByID(1):\ngot: %#v\nwant: %#v", got2, got1)
	}
	if mock.getByIDCalls != 1 {
		t.Errorf("GetByID() calls = %d, want %d", mock.getByIDCalls, 1)
	}

	_, err = c.GetByID(ctx, 999)
	if err == nil {
		t.Errorf("GetByID(999) error = %v, wantErr %v", err, true)
	}
	if mock.getByIDCalls != 2 {
		t.Errorf("GetByID() calls = %d, want %d", mock.getByIDCalls, 2)
	}
}

func TestModelCache_TTL_Expiration_GetByID(t *testing.T) {
	ctx := context.Background()
	mock := &mockProvider{
		models: map[int]domain.NNModel{
			1: {ID: 1, Name: "ModelA"},
		},
	}

	c := cache.NewModelCache(mock, 1*time.Millisecond)

	_, err := c.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID(1): unexpected error: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err = c.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID(1): unexpected error: %v", err)
	}

	if mock.getByIDCalls != 2 {
		t.Errorf("GetByID() calls = %d, want %d", mock.getByIDCalls, 2)
	}
}
