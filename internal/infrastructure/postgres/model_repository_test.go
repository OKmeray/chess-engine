//go:build integration

package postgres

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	testcontainerspg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB initializes a PostgreSQL container, applies migrations,
// and returns a configured ModelRepository and a teardown closure.
func setupTestDB(t *testing.T) (*ModelRepository, func()) {
	ctx := context.Background()

	pgContainer, err := testcontainerspg.Run(ctx,
		"postgres:17-alpine",
		testcontainerspg.WithDatabase("testdb"),
		testcontainerspg.WithUsername("testuser"),
		testcontainerspg.WithPassword("testpass"),
		testcontainers.WithLogger(tclog.TestLogger(t)),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to build connection string: %v", err)
	}

	db, err := New(connStr, 5, 5, 5*time.Minute)
	if err != nil {
		t.Fatalf("failed to open database connection pool: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("failed to initialize migration driver: %v", err)
	}

	migrationsPath, err := filepath.Abs("../../../db/migrations")
	if err != nil {
		t.Fatalf("failed to resolve migrations absolute path: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.ToSlash(migrationsPath),
		"postgres", driver,
	)
	if err != nil {
		t.Fatalf("failed to configure migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to execute up migrations: %v", err)
	}

	teardown := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return NewModelRepository(db), teardown
}

func TestModelRepository_GetByID(t *testing.T) {
	repo, teardown := setupTestDB(t)
	defer teardown()

	ctx := context.Background()
	testModelName := "integration-test-model"
	testModelPath := "models/test.onnx"

	dummyDetails := domain.NNModelDetails{
		Name:         testModelName,
		Architecture: "CNN",
		CNN: &domain.CNNDetails{
			Filters:   64,
			ResBlocks: 6,
		},
	}

	var insertedID int
	err := repo.db.QueryRowContext(ctx, `
		INSERT INTO models (name, path, details) 
		VALUES ($1, $2, $3) RETURNING id`,
		testModelName, testModelPath, dummyDetails,
	).Scan(&insertedID)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	gotModel, err := repo.GetByID(ctx, insertedID)
	if err != nil {
		t.Fatalf("GetByID(%d): unexpected error: %v", insertedID, err)
	}

	if gotModel.Path != testModelPath {
		t.Errorf("GetByID().Path = %q, want %q", gotModel.Path, testModelPath)
	}

	if gotModel.Details.Name != dummyDetails.Name {
		t.Errorf("GetByID().Details.Name = %q, want %q", gotModel.Details.Name, dummyDetails.Name)
	}

	if gotModel.Details.CNN == nil || gotModel.Details.CNN.Filters != dummyDetails.CNN.Filters {
		t.Errorf("GetByID().Details.CNN:\ngot:  %#v\nwant: %#v", gotModel.Details.CNN, dummyDetails.CNN)
	}
}

func TestModelRepository_GetByID_NotFound(t *testing.T) {
	repo, teardown := setupTestDB(t)
	defer teardown()

	_, err := repo.GetByID(context.Background(), 99999)
	if err == nil {
		t.Errorf("GetByID(%d) error = %v, wantErr %v", 99999, err, true)
	}
}

func TestModelRepository_GetAll(t *testing.T) {
	repo, teardown := setupTestDB(t)
	defer teardown()

	ctx := context.Background()

	_, err := repo.db.ExecContext(ctx, "TRUNCATE TABLE models RESTART IDENTITY")
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	dummyJSON := `{"architecture": "CNN", "cnn": {"filters": 64}}`
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO models (name, path, details) VALUES 
		('model_a', 'models/a.onnx', $1),
		('model_b', 'models/b.onnx', $2)`,
		dummyJSON, dummyJSON,
	)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	models, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll(): unexpected error: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Len: got %d; want %d", len(models), 2)
	}

	if len(models) > 0 {
		m := models[0]
		if m.Name != "model_a" && m.Name != "model_b" {
			t.Errorf("GetAll()[0].Name = %q, want %q or %q", m.Name, "model_a", "model_b")
		}
	}
}

func TestPostgres_New_FailsOnBadConnectionString(t *testing.T) {
	invalidConn := "invalid-conn-string"
	_, err := New(invalidConn, 1, 1, time.Second)
	if err == nil {
		t.Errorf("New(%q) error = %v, wantErr %v", invalidConn, err, true)
	}
}
