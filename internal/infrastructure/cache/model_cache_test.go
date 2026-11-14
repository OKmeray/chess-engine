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
	getByIDCalls   int
	getActiveCalls int
	models         map[int]domain.NNModel
	activeModel    domain.NNModel
	err            error
}

func (m *mockProvider) GetByID(ctx context.Context, id int) (domain.NNModel, error) {
	m.getByIDCalls++
	if m.err != nil {
		return domain.NNModel{}, m.err
	}
	model, ok := m.models[id]
	if !ok {
		return domain.NNModel{}, domain.ErrModelNotFound
	}
	return model, nil
}

func (m *mockProvider) GetActive(ctx context.Context) (domain.NNModel, error) {
	m.getActiveCalls++
	if m.err != nil {
		return domain.NNModel{}, m.err
	}
	if m.activeModel.ID != 0 {
		return m.activeModel, nil
	}
	return domain.NNModel{}, domain.ErrModelNotFound
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

func TestModelCache_GetActive(t *testing.T) {
	ctx := context.Background()
	activeModel := domain.NNModel{ID: 42, Name: "ActiveModel", IsActive: true}
	mock := &mockProvider{
		activeModel: activeModel,
		models: map[int]domain.NNModel{
			42: activeModel,
		},
	}

	c := cache.NewModelCache(mock, 1*time.Minute)

	// First call to GetActive hits provider
	got1, err := c.GetActive(ctx)
	if err != nil {
		t.Fatalf("GetActive(): unexpected error: %v", err)
	}
	if got1.ID != activeModel.ID {
		t.Errorf("GetActive().ID = %d, want %d", got1.ID, activeModel.ID)
	}
	if mock.getActiveCalls != 1 {
		t.Errorf("GetActive() calls = %d, want %d", mock.getActiveCalls, 1)
	}

	// Second call hits cache
	got2, err := c.GetActive(ctx)
	if err != nil {
		t.Fatalf("GetActive(): unexpected error: %v", err)
	}
	if got2.ID != activeModel.ID {
		t.Errorf("GetActive().ID = %d, want %d", got2.ID, activeModel.ID)
	}
	if mock.getActiveCalls != 1 {
		t.Errorf("GetActive() calls = %d, want %d", mock.getActiveCalls, 1)
	}

	// Calling GetByID with ActiveModelSentinelID should also return the cached active model
	got3, err := c.GetByID(ctx, domain.ActiveModelSentinelID)
	if err != nil {
		t.Fatalf("GetByID(%d): unexpected error: %v", domain.ActiveModelSentinelID, err)
	}
	if got3.ID != activeModel.ID {
		t.Errorf("GetByID(%d).ID = %d, want %d", domain.ActiveModelSentinelID, got3.ID, activeModel.ID)
	}
	if mock.getActiveCalls != 1 {
		t.Errorf("GetActive() calls = %d, want %d", mock.getActiveCalls, 1)
	}

	// Provider error returns error
	mockErr := &mockProvider{err: domain.ErrModelNotFound}
	cErr := cache.NewModelCache(mockErr, 1*time.Minute)
	_, err = cErr.GetActive(ctx)
	if err == nil {
		t.Errorf("GetActive() error = %v, wantErr %v", err, true)
	}
}
