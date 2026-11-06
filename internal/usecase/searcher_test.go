package usecase

import (
	"context"
	"errors"
	"testing"
)

type mockEvaluatorProvider struct {
	evaluator Evaluator
	err       error
}

func (m *mockEvaluatorProvider) Get(_ context.Context, _ int) (Evaluator, error) {
	return m.evaluator, m.err
}

func TestGameCreator_CreateGame_MCTS(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockEval := &mockEvaluator{}
		provider := &mockEvaluatorProvider{evaluator: mockEval}
		creator := NewGameCreator(provider, 1000, 2, 4)

		game, err := creator.CreateGame(ctx, EngineTypeMCTS, 1)
		if err != nil {
			t.Fatalf("CreateGame(%q): unexpected error: %v", EngineTypeMCTS, err)
		}
		if game == nil {
			t.Fatalf("CreateGame(%q) = nil, want non-nil MoveFinder", EngineTypeMCTS)
		}
	})

	t.Run("ProviderFailure", func(t *testing.T) {
		wantErr := errors.New("provider failure")
		errProvider := &mockEvaluatorProvider{err: wantErr}
		errCreator := NewGameCreator(errProvider, 1000, 2, 4)

		_, err := errCreator.CreateGame(ctx, EngineTypeMCTS, 1)
		if !errors.Is(err, wantErr) {
			t.Fatalf("CreateGame(%q) error = %v, wantErr %v", EngineTypeMCTS, err, wantErr)
		}
	})

	t.Run("NilProvider", func(t *testing.T) {
		nilCreator := NewGameCreator(nil, 1000, 2, 4)

		_, err := nilCreator.CreateGame(ctx, EngineTypeMCTS, 1)
		if !errors.Is(err, ErrNilEvaluatorProvider) {
			t.Fatalf("CreateGame(%q) error = %v, wantErr %v", EngineTypeMCTS, err, ErrNilEvaluatorProvider)
		}
	})
}

func TestGameCreator_CreateGame_Minimax(t *testing.T) {
	creator := NewGameCreator(&mockEvaluatorProvider{}, 1000, 2, 4)

	ctx := context.Background()
	game, err := creator.CreateGame(ctx, EngineTypeMinimax, 0)
	if err != nil {
		t.Fatalf("CreateGame(%q): unexpected error: %v", EngineTypeMinimax, err)
	}
	if game == nil {
		t.Fatalf("CreateGame(%q) = nil, want non-nil MoveFinder", EngineTypeMinimax)
	}
}

func TestGameCreator_CreateGame_Unsupported(t *testing.T) {
	creator := NewGameCreator(&mockEvaluatorProvider{}, 1000, 2, 4)

	ctx := context.Background()
	_, err := creator.CreateGame(ctx, EngineType("unsupported_type"), 0)
	if !errors.Is(err, ErrUnsupportedEngineType) {
		t.Fatalf("CreateGame(%q) error = %v, wantErr %v", "unsupported_type", err, ErrUnsupportedEngineType)
	}
}
