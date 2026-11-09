package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

type mockSearcher struct {
	move        engine.Move
	err         error
	calls       int
	recordedCtx context.Context
}

func (m *mockSearcher) Search(ctx context.Context, _ *engine.Position) (engine.Move, error) {
	m.calls++
	m.recordedCtx = ctx
	return m.move, m.err
}

func TestCalculateTargetTime(t *testing.T) {
	tests := []struct {
		name          string
		turn          int
		moves         int
		remainingTime time.Duration
		increment     time.Duration
		want          time.Duration
	}{
		{
			name:          "Turn 10, 30s remaining, 30 moves, no increment",
			turn:          10,
			moves:         30,
			remainingTime: 30 * time.Second,
			increment:     0,
			want:          1000 * time.Millisecond, // 30000 / (40 - 10) = 1000ms
		},
		{
			name:          "Turn 10, 30s remaining, 30 moves, 2s increment",
			turn:          10,
			moves:         30,
			remainingTime: 30 * time.Second,
			increment:     2 * time.Second,
			want:          2500 * time.Millisecond, // 1000ms + (2000 * 3 / 4) = 2500ms
		},
		{
			name:          "Complex position (> 45 legal moves)",
			turn:          10,
			moves:         50,
			remainingTime: 30 * time.Second,
			increment:     0,
			want:          1250 * time.Millisecond, // 1000 * 5 / 4 = 1250ms
		},
		{
			name:          "Simple position (< 15 legal moves)",
			turn:          10,
			moves:         10,
			remainingTime: 30 * time.Second,
			increment:     0,
			want:          750 * time.Millisecond, // 1000 * 3 / 4 = 750ms
		},
		{
			name:          "Endgame (Turn 35), 10s remaining (min divisor 20)",
			turn:          35,
			moves:         20,
			remainingTime: 10 * time.Second,
			increment:     0,
			want:          500 * time.Millisecond, // 10000 / 20 = 500ms
		},
		{
			name:          "Safety margin cap (near zero time)",
			turn:          20,
			moves:         20,
			remainingTime: 250 * time.Millisecond,
			increment:     0,
			want:          100 * time.Millisecond,
		},
		{
			name:          "Hard limit below minSearchTime with remainingTime above safetyMargin",
			turn:          20,
			moves:         20,
			remainingTime: 550 * time.Millisecond,
			increment:     0,
			want:          100 * time.Millisecond,
		},
		{
			name:          "Extremely low time clamp (prevent flagging)",
			turn:          20,
			moves:         20,
			remainingTime: 50 * time.Millisecond,
			increment:     0,
			want:          25 * time.Millisecond, // 50 / 2 = 25ms, preserves remaining clock
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTargetTime(tt.turn, tt.moves, tt.remainingTime, tt.increment)
			if got != tt.want {
				t.Errorf("calculateTargetTime(%d, %d, %v, %v) = %v, want %v",
					tt.turn, tt.moves, tt.remainingTime, tt.increment, got, tt.want)
			}
		})
	}
}

func TestNewGame_NilSearcher(t *testing.T) {
	game, err := NewGame(nil)
	if !errors.Is(err, ErrNilSearcher) {
		t.Fatalf("NewGame(nil) error = %v, wantErr %v", err, ErrNilSearcher)
	}
	if game != nil {
		t.Errorf("NewGame(nil) = %v, want %v", game, nil)
	}
}

func TestGame_FindBestMove(t *testing.T) {
	expectedMove := engine.NewMove(52, 36, engine.FlagDoublePawn) // e2e4
	mock := &mockSearcher{move: expectedMove}

	g, err := NewGame(mock)
	if err != nil {
		t.Fatalf("NewGame(): unexpected error: %v", err)
	}
	ctx := context.Background()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	remainingTime := 60 * time.Second
	increment := 2 * time.Second
	startTime := time.Now()

	move, err := g.FindBestMove(ctx, fen, remainingTime, increment)
	if err != nil {
		t.Fatalf("FindBestMove(%q): unexpected error: %v", fen, err)
	}

	if mock.calls != 1 {
		t.Errorf("mockSearcher.calls: got %d, want %d", mock.calls, 1)
	}

	want := expectedMove.UCI()
	if move != want {
		t.Errorf("FindBestMove(%q) = %q, want %q", fen, move, want)
	}

	// Verify deadline invariants deterministically
	if mock.recordedCtx == nil {
		t.Fatalf("Search() context was not recorded")
	}
	deadline, ok := mock.recordedCtx.Deadline()
	if !ok {
		t.Errorf("Search() context has no deadline, want deadline set")
	} else {
		if !deadline.After(startTime) {
			t.Errorf("Search() context deadline %v is not after start %v", deadline, startTime)
		}
		if deadline.After(startTime.Add(remainingTime)) {
			t.Errorf("Search() context deadline %v exceeds remaining time %v", deadline, remainingTime)
		}
	}

	// Verify context cleanup (proves 'defer cancel()' ran and freed timer goroutines)
	if !errors.Is(mock.recordedCtx.Err(), context.Canceled) {
		t.Errorf("Search() context error = %v, wantErr %v", mock.recordedCtx.Err(), context.Canceled)
	}
}

func TestGame_FindBestMove_EmptyMoveReturnsError(t *testing.T) {
	mock := &mockSearcher{move: engine.Move(0)}

	g, err := NewGame(mock)
	if err != nil {
		t.Fatalf("NewGame(): unexpected error: %v", err)
	}
	ctx := context.Background()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	move, err := g.FindBestMove(ctx, fen, 60*time.Second, 0)
	if !errors.Is(err, ErrSearchFailed) {
		t.Fatalf("FindBestMove(%q) error = %v, wantErr %v", fen, err, ErrSearchFailed)
	}
	if move != "" {
		t.Errorf("FindBestMove(%q) move = %q, want %q", fen, move, "")
	}

	if mock.calls != 1 {
		t.Errorf("mockSearcher.calls: got %d, want %d", mock.calls, 1)
	}
}

func TestGame_FindBestMove_IllegalMoveReturnsError(t *testing.T) {
	// An illegal move for the starting position (h1 to e4 is impossible)
	illegalMove := engine.NewMove(63, 36, 0)
	mock := &mockSearcher{move: illegalMove}

	g, err := NewGame(mock)
	if err != nil {
		t.Fatalf("NewGame(): unexpected error: %v", err)
	}
	ctx := context.Background()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	move, err := g.FindBestMove(ctx, fen, 60*time.Second, 0)
	if !errors.Is(err, ErrSearchFailed) {
		t.Fatalf("FindBestMove(%q) error = %v, wantErr %v", fen, err, ErrSearchFailed)
	}
	if move != "" {
		t.Errorf("FindBestMove(%q) move = %q, want %q", fen, move, "")
	}

	if mock.calls != 1 {
		t.Errorf("mockSearcher.calls: got %d, want %d", mock.calls, 1)
	}
}

func TestGame_FindBestMove_OnlyOneLegalMove(t *testing.T) {
	// White king on h1 in check by black queen on g2. Only legal move is h1xg2.
	fen := "k7/8/8/8/8/8/6q1/7K w - - 0 1"
	mock := &mockSearcher{}

	g, err := NewGame(mock)
	if err != nil {
		t.Fatalf("NewGame(): unexpected error: %v", err)
	}
	ctx := context.Background()

	moveUCI, err := g.FindBestMove(ctx, fen, 60*time.Second, 0)
	if err != nil {
		t.Fatalf("FindBestMove(%q): unexpected error: %v", fen, err)
	}

	// Should not have called searcher at all.
	if mock.calls != 0 {
		t.Errorf("mockSearcher.calls: got %d, want %d", mock.calls, 0)
	}

	want := "h1g2"
	if moveUCI != want {
		t.Errorf("FindBestMove(%q) = %q, want %q", fen, moveUCI, want)
	}
}

func TestGame_FindBestMove_Errors(t *testing.T) {
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Already cancelled

	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	stalemateFEN := "7k/5K2/6Q1/8/8/8/8/8 b - - 0 80"

	tests := []struct {
		name      string
		ctx       context.Context
		searcher  Searcher
		fen       string
		wantErr   error
		wantCalls int
	}{
		{
			name:      "Nil searcher",
			ctx:       context.Background(),
			searcher:  nil,
			fen:       startFEN,
			wantErr:   ErrNilSearcher,
			wantCalls: 0,
		},
		{
			name:      "Pre-cancelled context",
			ctx:       cancelledCtx,
			searcher:  &mockSearcher{},
			fen:       startFEN,
			wantErr:   context.Canceled,
			wantCalls: 0,
		},
		{
			name:      "Invalid FEN string",
			ctx:       context.Background(),
			searcher:  &mockSearcher{},
			fen:       "invalid_fen_string",
			wantErr:   ErrInvalidFEN,
			wantCalls: 0,
		},
		{
			name:      "No legal moves (stalemate)",
			ctx:       context.Background(),
			searcher:  &mockSearcher{},
			fen:       stalemateFEN,
			wantErr:   ErrNoLegalMoves,
			wantCalls: 0,
		},
		{
			name:      "Searcher failure",
			ctx:       context.Background(),
			searcher:  &mockSearcher{err: errors.New("search failed")},
			fen:       startFEN,
			wantErr:   ErrSearchFailed,
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g *Game
			if tt.searcher != nil {
				var err error
				g, err = NewGame(tt.searcher)
				if err != nil {
					t.Fatalf("NewGame(): unexpected error: %v", err)
				}
			} else {
				g = &Game{}
			}

			_, err := g.FindBestMove(tt.ctx, tt.fen, 60*time.Second, 0)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("FindBestMove(%q) error = %v, wantErr %v", tt.fen, err, tt.wantErr)
			}

			if m, ok := tt.searcher.(*mockSearcher); ok {
				if m.calls != tt.wantCalls {
					t.Errorf("mockSearcher.calls: got %d, want %d", m.calls, tt.wantCalls)
				}
			}
		})
	}
}
