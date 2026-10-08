package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/OKmeray/chess-engine/internal/domain"
)

// ModelRepository provides a PostgreSQL-backed implementation for managing models.
type ModelRepository struct {
	db *sql.DB
}

// NewModelRepository initializes the repository with a database connection pool.
func NewModelRepository(db *sql.DB) *ModelRepository {
	return &ModelRepository{
		db: db,
	}
}

// GetByID returns a model by its ID.
func (r *ModelRepository) GetByID(ctx context.Context, id int) (domain.NNModel, error) {
	query := `SELECT id, name, path, details, is_active, created_at FROM models WHERE id = $1`

	var model domain.NNModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ID,
		&model.Name,
		&model.Path,
		&model.Details,
		&model.IsActive,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NNModel{}, fmt.Errorf("%w: id %d", domain.ErrModelNotFound, id)
		}
		return domain.NNModel{}, fmt.Errorf("error querying model: %w", err)
	}

	return model, nil
}

// GetAll returns all available models.
func (r *ModelRepository) GetAll(ctx context.Context) ([]domain.NNModel, error) {
	query := `SELECT id, name, path, details, is_active, created_at FROM models`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying models: %w", err)
	}
	defer rows.Close()

	var models []domain.NNModel

	for rows.Next() {
		var m domain.NNModel
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Path,
			&m.Details,
			&m.IsActive,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning model row: %w", err)
		}
		models = append(models, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating model rows: %w", err)
	}

	return models, nil
}
