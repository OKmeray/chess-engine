package onnx

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/OKmeray/chess-engine/internal/domain"
)

type mockFetcher struct {
	mu        sync.Mutex
	models    map[int]domain.NNModel
	err       error
	callCount int
}

func (m *mockFetcher) GetByID(_ context.Context, id int) (domain.NNModel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.err != nil {
		return domain.NNModel{}, m.err
	}
	model, ok := m.models[id]
	if !ok {
		return domain.NNModel{}, domain.ErrModelNotFound
	}
	return model, nil
}

func (m *mockFetcher) calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestManager_Get_SuccessAndCache(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	modelID := 1
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID: {ID: modelID, Name: "model_1", Path: modelPath},
		},
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	ctx := context.Background()

	// Initial fetch
	eval1, err := mgr.Get(ctx, modelID)
	if err != nil {
		t.Fatalf("Get(%d): unexpected error: %v", modelID, err)
	}
	if eval1 == nil {
		t.Fatalf("Get(%d) eval = nil, want non-nil", modelID)
	}

	// Cached fetch
	eval2, err := mgr.Get(ctx, modelID)
	if err != nil {
		t.Fatalf("Get(%d) cached: unexpected error: %v", modelID, err)
	}

	if eval1 != eval2 {
		t.Errorf("Get(%d) instance mismatch: got %p, want %p", modelID, eval2, eval1)
	}

	if gotCalls := fetcher.calls(); gotCalls != 1 {
		t.Errorf("fetcher.calls() = %d, want %d", gotCalls, 1)
	}
}

func TestManager_Get_CanceledContext(t *testing.T) {
	modelID := 1
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID: {ID: modelID, Name: "model_1", Path: "model.onnx"},
		},
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := mgr.Get(ctx, modelID)
	if err == nil {
		t.Fatalf("Get(%d) with canceled ctx error = nil, wantErr %v", modelID, true)
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Get(%d) error = %v, wantErr %v", modelID, err, context.Canceled)
	}

	if gotCalls := fetcher.calls(); gotCalls != 0 {
		t.Errorf("fetcher.calls() = %d, want %d", gotCalls, 0)
	}
}

func TestManager_Get_FetcherError(t *testing.T) {
	modelID := 999
	fetcher := &mockFetcher{
		err: domain.ErrModelNotFound,
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	_, err := mgr.Get(context.Background(), modelID)
	if err == nil {
		t.Fatalf("Get(%d) error = nil, wantErr %v", modelID, true)
	}
	if !errors.Is(err, ErrFetchModelMetadata) {
		t.Errorf("Get(%d) error = %v, wantErr %v", modelID, err, ErrFetchModelMetadata)
	}
	if !errors.Is(err, domain.ErrModelNotFound) {
		t.Errorf("Get(%d) wrapped error = %v, wantErr %v", modelID, err, domain.ErrModelNotFound)
	}
}

func TestManager_Get_BadONNXFile(t *testing.T) {
	initTestEnvironment(t)

	modelID := 2
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID: {ID: modelID, Name: "bad_model", Path: "non_existent_file.onnx"},
		},
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	_, err := mgr.Get(context.Background(), modelID)
	if err == nil {
		t.Fatalf("Get(%d) error = nil, wantErr %v", modelID, true)
	}
	if !errors.Is(err, ErrLoadONNXFile) {
		t.Errorf("Get(%d) error = %v, wantErr %v", modelID, err, ErrLoadONNXFile)
	}
}

func TestManager_Get_Concurrent_Singleflight(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	modelID := 1
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID: {ID: modelID, Name: "model_1", Path: modelPath},
		},
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	const concurrency = 100
	var wg sync.WaitGroup
	errCh := make(chan error, concurrency)
	startGate := make(chan struct{})

	for range concurrency {
		wg.Go(func() {
			<-startGate
			_, err := mgr.Get(context.Background(), modelID)
			if err != nil {
				errCh <- err
			}
		})
	}

	close(startGate)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("concurrent Get(%d): unexpected error: %v", modelID, err)
	}

	if gotCalls := fetcher.calls(); gotCalls != 1 {
		t.Errorf("fetcher.calls() = %d, want %d", gotCalls, 1)
	}
}

func TestManager_Get_Concurrent_DifferentModels(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	modelID1, modelID2 := 1, 2
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID1: {ID: modelID1, Name: "model_1", Path: modelPath},
			modelID2: {ID: modelID2, Name: "model_2", Path: modelPath},
		},
	}

	mgr := NewManager(fetcher)
	defer mgr.CloseAll()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	startGate := make(chan struct{})

	for _, id := range []int{modelID1, modelID2} {
		modelID := id
		wg.Go(func() {
			<-startGate
			_, err := mgr.Get(context.Background(), modelID)
			if err != nil {
				errCh <- err
			}
		})
	}

	close(startGate)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("concurrent multi-model Get: unexpected error: %v", err)
	}

	if gotCalls := fetcher.calls(); gotCalls != 2 {
		t.Errorf("fetcher.calls() = %d, want %d", gotCalls, 2)
	}
}

func TestManager_CloseAll(t *testing.T) {
	initTestEnvironment(t)
	modelPath := getTestModelPath(t)

	modelID := 1
	fetcher := &mockFetcher{
		models: map[int]domain.NNModel{
			modelID: {ID: modelID, Name: "model_1", Path: modelPath},
		},
	}

	mgr := NewManager(fetcher)

	ctx := context.Background()
	_, err := mgr.Get(ctx, modelID)
	if err != nil {
		t.Fatalf("Get(%d): unexpected error: %v", modelID, err)
	}

	if gotLen := len(mgr.loaded); gotLen != 1 {
		t.Errorf("Len: got %d, want %d", gotLen, 1)
	}

	mgr.CloseAll()

	if gotLen := len(mgr.loaded); gotLen != 0 {
		t.Errorf("Len: got %d, want %d", gotLen, 0)
	}
}
