package cache_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain"
	"github.com/OKmeray/chess-engine/internal/infrastructure/cache"
)

type mockProvider struct {
	getByIDCalls int
	getAllCalls  int
	getAllErr    error
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

func (m *mockProvider) GetAll(ctx context.Context) ([]domain.NNModel, error) {
	m.getAllCalls++
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	var all []domain.NNModel
	for _, v := range m.models {
		all = append(all, v)
	}
	return all, nil
}

func TestModelCache_GetAll(t *testing.T) {
	ctx := context.Background()
	mock := &mockProvider{
		models: map[int]domain.NNModel{
			1: {ID: 1, Name: "ModelA"},
			2: {ID: 2, Name: "ModelB"},
		},
	}

	c := cache.NewModelCache(mock, 1*time.Minute)

	models1, err := c.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll(): unexpected error: %v", err)
	}
	if len(models1) != 2 {
		t.Errorf("Len: got %d, want %d", len(models1), 2)
	}
	if mock.getAllCalls != 1 {
		t.Errorf("GetAll() calls = %d, want %d", mock.getAllCalls, 1)
	}

	models2, err := c.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll(): unexpected error: %v", err)
	}
	if len(models2) != 2 {
		t.Errorf("Len: got %d, want %d", len(models2), 2)
	}
	if mock.getAllCalls != 1 {
		t.Errorf("GetAll() calls = %d, want %d", mock.getAllCalls, 1)
	}
}

func TestModelCache_GetAll_Error(t *testing.T) {
	ctx := context.Background()
	mock := &mockProvider{
		getAllErr: errors.New("database down"),
	}

	c := cache.NewModelCache(mock, 1*time.Minute)

	_, err := c.GetAll(ctx)
	if err == nil {
		t.Errorf("GetAll() error = %v, wantErr %v", err, true)
	}
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

func TestModelCache_TTL_Expiration_GetAll(t *testing.T) {
	ctx := context.Background()
	mock := &mockProvider{
		models: map[int]domain.NNModel{
			1: {ID: 1, Name: "ModelA"},
		},
	}

	c := cache.NewModelCache(mock, 1*time.Millisecond)

	_, err := c.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll(): unexpected error: %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, err = c.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll(): unexpected error: %v", err)
	}

	if mock.getAllCalls != 2 {
		t.Errorf("GetAll() calls = %d, want %d", mock.getAllCalls, 2)
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
