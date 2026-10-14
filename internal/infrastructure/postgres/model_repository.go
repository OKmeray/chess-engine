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

// List returns a paginated list of models and the total count.
func (r *ModelRepository) List(ctx context.Context, limit, offset int) (models []domain.NNModel, total int, err error) {
	countQuery := `SELECT COUNT(*) FROM models`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("error counting models: %w", err)
	}

	query := `SELECT id, name, path, details, is_active, created_at FROM models ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying models: %w", err)
	}
	defer rows.Close()

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
			return nil, 0, fmt.Errorf("error scanning model row: %w", err)
		}
		models = append(models, m)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating model rows: %w", err)
	}

	return models, total, nil
}
