package onnx

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/OKmeray/chess-engine/internal/domain"
	"github.com/OKmeray/chess-engine/internal/usecase"
	"golang.org/x/sync/singleflight"
)

// ModelFetcher defines how the manager fetches model metadata.
type ModelFetcher interface {
	GetByID(ctx context.Context, id int) (domain.NNModel, error)
}

// Manager handles the lazy-loading and caching of ONNX Evaluators.
type Manager struct {
	mu      sync.RWMutex
	fetcher ModelFetcher
	loaded  map[int]*Evaluator
	group   singleflight.Group
}

// NewManager creates a new Manager.
func NewManager(fetcher ModelFetcher) *Manager {
	return &Manager{
		fetcher: fetcher,
		loaded:  make(map[int]*Evaluator),
	}
}

// Get returns a live ONNX Evaluator for the given model ID.
func (m *Manager) Get(ctx context.Context, modelID int) (usecase.Evaluator, error) {
	m.mu.RLock()
	eval, exists := m.loaded[modelID]
	m.mu.RUnlock()

	if exists {
		return eval, nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	key := strconv.Itoa(modelID)
	val, err, _ := m.group.Do(key, func() (any, error) {
		model, err := m.fetcher.GetByID(ctx, modelID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFetchModelMetadata, err)
		}

		newEval, err := NewEvaluator(model.Path)
		if err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrLoadONNXFile, model.Path, err)
		}

		m.mu.Lock()
		m.loaded[modelID] = newEval
		m.mu.Unlock()

		return newEval, nil
	})
	if err != nil {
		return nil, err
	}

	evaluator, ok := val.(usecase.Evaluator)
	if !ok {
		return nil, fmt.Errorf("unexpected evaluator type: %T", val)
	}

	return evaluator, nil
}

// CloseAll safely destroys all loaded ONNX sessions.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, eval := range m.loaded {
		eval.Close()
		delete(m.loaded, id)
	}
}
