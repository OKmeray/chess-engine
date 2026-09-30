package cache

import (
	"context"
	"sync"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain"
)

// ModelProvider defines the data access contract required by the cache.
type ModelProvider interface {
	GetByID(ctx context.Context, id int) (domain.NNModel, error)
	GetAll(ctx context.Context) ([]domain.NNModel, error)
}

type cacheItem struct {
	model     domain.NNModel
	fetchedAt time.Time
}

// ModelCache provides a thread-safe, TTL-based memory cache decorating a ModelProvider.
type ModelCache struct {
	provider ModelProvider
	ttl      time.Duration

	mu           sync.RWMutex
	models       map[int]cacheItem
	all          []domain.NNModel
	allFetchedAt time.Time
}

// NewModelCache initializes a new ModelCache wrapper.
func NewModelCache(provider ModelProvider, ttl time.Duration) *ModelCache {
	return &ModelCache{
		provider: provider,
		ttl:      ttl,
		models:   make(map[int]cacheItem),
	}
}

// GetAll returns all models.
func (c *ModelCache) GetAll(ctx context.Context) ([]domain.NNModel, error) {
	c.mu.RLock()
	if len(c.all) > 0 && time.Since(c.allFetchedAt) < c.ttl {
		res := c.all
		c.mu.RUnlock()
		return res, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check condition after acquiring write lock
	if len(c.all) > 0 && time.Since(c.allFetchedAt) < c.ttl {
		return c.all, nil
	}

	models, err := c.provider.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	c.all = models
	c.allFetchedAt = now

	for _, m := range models {
		c.models[m.ID] = cacheItem{
			model:     m,
			fetchedAt: now,
		}
	}

	return c.all, nil
}

// GetByID returns a model by its ID.
func (c *ModelCache) GetByID(ctx context.Context, id int) (domain.NNModel, error) {
	c.mu.RLock()
	item, exists := c.models[id]
	if exists && time.Since(item.fetchedAt) < c.ttl {
		c.mu.RUnlock()
		return item.model, nil
	}
	c.mu.RUnlock()

	model, err := c.provider.GetByID(ctx, id)
	if err != nil {
		return domain.NNModel{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.models[id] = cacheItem{
		model:     model,
		fetchedAt: time.Now(),
	}

	return model, nil
}
