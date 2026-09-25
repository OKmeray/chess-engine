package usecase

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

func TestMCTSNode_AddValue(t *testing.T) {
	node := &MCTSNode{}
	var wg sync.WaitGroup

	numGoroutines := 100
	for i := 0; i < numGoroutines; i++ {
		wg.Go(func() {
			node.addValue(1.5)
		})
	}
	wg.Wait()

	got := node.getValue()
	want := float32(150.0)
	if math.Abs(float64(got-want)) > 1e-6 {
		t.Errorf("MCTSNode.getValue() = %f, want %f", got, want)
	}
}

func TestMCTSSearcher_Backpropagate(t *testing.T) {
	m := NewMCTSSearcher(nil, 10, 1)
	// history: root(0), child(1), grandchild(2), great-grandchild(3)
	history := []int32{0, 1, 2, 3}

	m.backpropagate(history, -1.0)

	// Great-Grandchild (3) should get 1.0
	if got := m.nodes[3].getValue(); got != 1.0 {
		t.Errorf("nodes[3].getValue() = %f, want 1.0", got)
	}
	// Grandchild (2) should get -1.0
	if got := m.nodes[2].getValue(); got != -1.0 {
		t.Errorf("nodes[2].getValue() = %f, want -1.0", got)
	}
	// Child (1) should get +1.0
	if got := m.nodes[1].getValue(); got != 1.0 {
		t.Errorf("nodes[1].getValue() = %f, want 1.0", got)
	}
	// Root (0) should get -1.0
	if got := m.nodes[0].getValue(); got != -1.0 {
		t.Errorf("nodes[0].getValue() = %f, want -1.0", got)
	}

	for i := 0; i < 4; i++ {
		if got := m.nodes[i].Visits; got != 1 {
			t.Errorf("nodes[%d].Visits = %d, want 1", i, got)
		}
	}
}

func TestMCTSSearcher_Expand(t *testing.T) {
	m := NewMCTSSearcher(nil, 10, 1)
	m.nodes[0] = MCTSNode{State: StateExpanding}
	m.nodeCount = 1

	// start pos
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	move1 := engine.NewMove(engine.GetNumBySquareName("e2"), engine.GetNumBySquareName("e4"), engine.FlagDoublePawn)
	move2 := engine.NewMove(engine.GetNumBySquareName("e2"), engine.GetNumBySquareName("e3"), engine.FlagQuiet)
	legalMoves := []engine.Move{move1, move2}

	idx1 := EncodeMove(move1, false)
	idx2 := EncodeMove(move2, false)

	resp := &EvalResponse{
		Policy: make([]float32, PolicySize),
		Value:  0.0,
	}
	resp.Policy[idx1] = 5.0
	resp.Policy[idx2] = 1.0

	m.expand(0, resp, pos.SideToMove, legalMoves)

	if got := m.nodes[0].NumChildren; got != int16(len(legalMoves)) {
		t.Errorf("nodes[0].NumChildren = %d, want %d", got, len(legalMoves))
	}

	child1 := m.nodes[1]
	child2 := m.nodes[2]

	sumPrior := child1.Prior + child2.Prior
	if math.Abs(float64(sumPrior-1.0)) > 1e-6 {
		t.Errorf("sum of priors = %f, want 1.0", sumPrior)
	}

	if child1.Prior <= child2.Prior {
		t.Errorf("child1.Prior = %f, want > child2.Prior (%f)", child1.Prior, child2.Prior)
	}

	if got := m.nodes[0].State; got != StateExpanded {
		t.Errorf("nodes[0].State = %d, want %d", got, StateExpanded)
	}
}

func TestMCTSSearcher_Expand_EmptyMoves(t *testing.T) {
	searcher := NewMCTSSearcher(nil, 1000000, 1)
	searcher.nodes[0] = MCTSNode{State: StateExpanding}

	searcher.expand(0, &EvalResponse{}, engine.White, []engine.Move{})

	if got := searcher.nodes[0].State; got != StateExpanded {
		t.Errorf("nodes[0].State = %v, want %v", got, StateExpanded)
	}
	if got := searcher.nodes[0].NumChildren; got != 0 {
		t.Errorf("nodes[0].NumChildren = %d, want 0", got)
	}
}

func TestMCTSSearcher_Expand_ArenaFull(t *testing.T) {
	// Searcher with capacity of 2 nodes
	searcher := NewMCTSSearcher(nil, 2, 1)

	searcher.nodes[0] = MCTSNode{State: StateExpanding}
	searcher.nodeCount = 1

	// Provide 2 legal moves. This requires allocating 2 nodes.
	// firstChildIdx would be 1. 1 + 2 = 3 > len(nodes)
	legalMoves := []engine.Move{engine.Move(1), engine.Move(2)}

	resp := &EvalResponse{
		Policy: make([]float32, PolicySize),
		Value:  0.0,
	}

	searcher.expand(0, resp, engine.White, legalMoves)

	if got := searcher.nodes[0].State; got != StateExpanded {
		t.Errorf("nodes[0].State = %v, want %v", got, StateExpanded)
	}
	if got := searcher.nodes[0].NumChildren; got != 0 {
		t.Errorf("nodes[0].NumChildren = %v, want 0", got)
	}
	if got := searcher.nodeCount; got != 1 {
		t.Errorf("searcher.nodeCount = %v, want 1", got)
	}
}

func TestMCTSSearcher_SelectChild(t *testing.T) {
	m := NewMCTSSearcher(nil, 10, 1)
	m.nodes[0] = MCTSNode{
		NumChildren:   3,
		FirstChildIdx: 1,
		Visits:        10,
	}

	// Child 1: High Prior, 0 Visits -> High U, 0 Q
	m.nodes[1] = MCTSNode{Prior: 0.9, Visits: 0}

	// Child 2: Low Prior, High Visits, High Value -> High Q, Low U
	m.nodes[2] = MCTSNode{Prior: 0.05, Visits: 5}
	m.nodes[2].addValue(4.5) // q = 4.5/5 = 0.9

	// Child 3: High Prior, but penalized by VirtualLoss
	m.nodes[3] = MCTSNode{Prior: 0.9, Visits: 0, VirtualLoss: 100}

	bestChild := m.selectChild(0)

	// Mathematically:
	// ParentN = 10. sqrt(10) = 3.16
	// child 1: q = 0. u = 0.9 * 8 * 3.16 / 1 = 22.7. score = 22.7
	// child 2: q = 0.9. u = 0.05 * 8 * 3.16 / 6 = 0.21. score = 1.11
	// child 3: totalVisits = 100. q = -100/100 = -1. u = 0.9 * 8 * 3.16 / 101 = 0.22. score = -0.78
	// Best score has child 1

	if bestChild != 1 {
		t.Errorf("selectChild(0) = %d, want 1", bestChild)
	}
}

type mockEvaluator struct {
	favoredMove engine.Move
}

func (e *mockEvaluator) EvaluateBatch(positions []engine.Position) ([][]float32, []float32, error) {
	policies := make([][]float32, len(positions))
	values := make([]float32, len(positions))

	for i, pos := range positions {
		policies[i] = make([]float32, PolicySize)
		idx := EncodeMove(e.favoredMove, pos.SideToMove == engine.Black)
		if idx >= 0 && idx < PolicySize {
			policies[i][idx] = 100.0 // Very high logit for the favored move
		}
		values[i] = 0.0
	}
	return policies, values, nil
}

func TestMCTSSearcher_Search(t *testing.T) {
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	favoredMove := engine.NewMove(52, 36, engine.FlagDoublePawn) // e2-e4

	evaluator := &mockEvaluator{
		favoredMove: favoredMove,
	}

	searcher := NewMCTSSearcher(evaluator, 1000000, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	bestMove, err := searcher.Search(ctx, pos)
	if err != nil {
		t.Fatalf("Search(): unexpected error: %v", err)
	}

	if bestMove != favoredMove {
		t.Errorf("Search() move = %v, want %v", bestMove, favoredMove)
	}

	if visits := atomic.LoadInt32(&searcher.nodes[0].Visits); visits <= 1 {
		t.Errorf("nodes[0].Visits = %d, want > 1", visits)
	}
}

func TestMCTSSearcher_Search_TerminalNodes(t *testing.T) {
	evaluator := &mockEvaluator{favoredMove: engine.NewMove(52, 36, engine.FlagDoublePawn)}

	tests := []struct {
		name          string
		fen           string
		expectedValue float32
	}{
		{
			name:          "Checkmate",
			fen:           "4k3/4Q3/4K3/8/8/8/8/8 b - - 0 101",
			expectedValue: 1.0,
		},
		{
			name:          "Stalemate",
			fen:           "k7/2Q5/8/8/8/8/8/7K b - - 0 101",
			expectedValue: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := engine.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			searcher := NewMCTSSearcher(evaluator, 1000000, 1)
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			_, err = searcher.Search(ctx, pos)
			if err != nil {
				t.Fatalf("Search(%v): unexpected error: %v", tt.name, err)
			}

			sum := searcher.nodes[0].getValue()
			visits := float32(atomic.LoadInt32(&searcher.nodes[0].Visits))
			got := sum / visits
			if got != tt.expectedValue {
				t.Errorf("average value = %v, want %v", got, tt.expectedValue)
			}
		})
	}
}

// TestMCTSSearcher_Search_DeepDraws verify that MCTS catches draws deep in the search tree (after a move is made)
func TestMCTSSearcher_Search_DeepDraws(t *testing.T) {
	evaluator := &mockEvaluator{favoredMove: engine.Move(0)}

	tests := []struct {
		name string
		fen  string
	}{
		{
			name: "Capture leads to Insufficient Material",
			fen:  "8/8/8/8/1K6/p7/8/1k6 w - - 1 103",
		},
		{
			name: "Any move leads to 50-Move Rule",
			fen:  "k6r/8/8/8/8/8/8/K7 w - - 99 105",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, _ := engine.ParseFEN(tt.fen)
			searcher := NewMCTSSearcher(evaluator, 1000000, 1)
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			searcher.Search(ctx, pos)

			// The root should have at least 1 child, and that child should be terminal (0.0)
			if searcher.nodes[0].NumChildren == 0 {
				t.Fatalf("nodes[0].NumChildren = 0, want > 0")
			}
			child := searcher.nodes[searcher.nodes[0].FirstChildIdx]

			sum := child.getValue()
			visits := float32(atomic.LoadInt32(&child.Visits))
			got := sum / visits

			if got != 0.0 {
				t.Errorf("getValue() / Visits = %v, want 0.0", got)
			}
		})
	}
}

var errMockEvaluator = errors.New("mock evaluator error")

type errorEvaluator struct{}

func (e *errorEvaluator) EvaluateBatch(positions []engine.Position) ([][]float32, []float32, error) {
	return nil, nil, errMockEvaluator
}

func TestMCTSSearcher_Search_EvaluatorError(t *testing.T) {
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	searcher := NewMCTSSearcher(&errorEvaluator{}, 1000000, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Verify that NN failures trigger an immediate hard abort and return the error.
	move, err := searcher.Search(ctx, pos)
	if !errors.Is(err, errMockEvaluator) {
		t.Errorf("Search() error = %v, wantErr %v", err, errMockEvaluator)
	}
	if move != engine.Move(0) {
		t.Errorf("Search() move = %v, want %v", move, engine.Move(0))
	}
}

func TestMCTSSearcher_Search_Contention(t *testing.T) {
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	favoredMove := engine.NewMove(52, 36, engine.FlagDoublePawn)
	// Use 50 workers to ensure high contention on StateUnexpanded -> StateExpanding
	searcher := NewMCTSSearcher(&mockEvaluator{favoredMove: favoredMove}, 1000000, 50)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	move, err := searcher.Search(ctx, pos)
	if err != nil {
		t.Fatalf("Search(pos): unexpected error: %v", err)
	}

	if move != favoredMove {
		t.Errorf("Search() move = %v, want %v", move, favoredMove)
	}

	if visits := atomic.LoadInt32(&searcher.nodes[0].Visits); visits < 50 {
		t.Errorf("nodes[0].Visits = %d, want >= 50", visits)
	}
}

type panicEvaluator struct{}

func (e *panicEvaluator) EvaluateBatch(positions []engine.Position) ([][]float32, []float32, error) {
	panic("simulated fatal panic in neural network batcher")
}

func TestMCTSSearcher_Search_PanicRecovery(t *testing.T) {
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	searcher := NewMCTSSearcher(&panicEvaluator{}, 1000000, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := searcher.Search(ctx, pos)
	if !errors.Is(err, ErrWorkerPanic) {
		t.Errorf("Search() error = %v, wantErr %v", err, ErrWorkerPanic)
	}
}

// slowEvaluator simulates a delayed neural network response to test timeout recovery.
type slowEvaluator struct {
	start chan struct{}
	ctx   context.Context
}

func (e *slowEvaluator) EvaluateBatch(positions []engine.Position) ([][]float32, []float32, error) {
	select {
	case <-e.start:
	default:
		close(e.start)
	}

	<-e.ctx.Done()

	// Allow tree workers to process the cancellation context before returning the error.
	time.Sleep(100 * time.Millisecond)
	return nil, nil, e.ctx.Err()
}

func TestMCTSSearcher_TreeWorker_ContextCancellation(t *testing.T) {
	pos, _ := engine.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eval := &slowEvaluator{start: make(chan struct{}), ctx: ctx}
	searcher := NewMCTSSearcher(eval, 1000000, 2)

	// Force channel contention by making the request queue unbuffered.
	searcher.evalReqCh = make(chan EvalRequest)

	// Pre-build a 3-node tree to prevent workers from colliding on the root node's expansion lock.
	atomic.StoreInt32(&searcher.nodeCount, 3)
	searcher.nodes[0] = MCTSNode{State: StateExpanded, FirstChildIdx: 1, NumChildren: 2}
	searcher.nodes[1] = MCTSNode{State: StateUnexpanded}
	searcher.nodes[2] = MCTSNode{State: StateUnexpanded}

	var wg sync.WaitGroup

	wg.Go(func() {
		searcher.nnBatchWorker(ctx)
	})

	var fatalErr atomic.Pointer[error]
	wg.Go(func() {
		searcher.treeWorker(ctx, pos, &fatalErr, cancel)
	})

	// Await the broadcast signal indicating the first worker is trapped inside the mock evaluator.
	<-eval.start

	wg.Go(func() {
		searcher.treeWorker(ctx, pos, &fatalErr, cancel)
	})

	// Allow the OS scheduler enough time to trap the second worker in the unbuffered channel queue.
	time.Sleep(50 * time.Millisecond)

	// Trigger cancellation to verify both blocked workers cleanly revert their VirtualLoss.
	cancel()

	wg.Wait()

	got0 := atomic.LoadInt32(&searcher.nodes[0].VirtualLoss)
	got1 := atomic.LoadInt32(&searcher.nodes[1].VirtualLoss)
	got2 := atomic.LoadInt32(&searcher.nodes[2].VirtualLoss)

	if got0 != 0 {
		t.Errorf("nodes[0].VirtualLoss = %d, want 0", got0)
	}
	if got1 != 0 {
		t.Errorf("nodes[1].VirtualLoss = %d, want 0", got1)
	}
	if got2 != 0 {
		t.Errorf("nodes[2].VirtualLoss = %d, want 0", got2)
	}
}

type mockEarlyExitEvaluator struct {
	firstCall bool
	mu        sync.Mutex
}

func (e *mockEarlyExitEvaluator) EvaluateBatch(positions []engine.Position) ([][]float32, []float32, error) {
	e.mu.Lock()
	if e.firstCall {
		e.firstCall = false
		e.mu.Unlock()
		// Delay to guarantee the 50ms ticker fires while root is unexpanded
		time.Sleep(100 * time.Millisecond)
	} else {
		e.mu.Unlock()
	}

	policies := make([][]float32, len(positions))
	values := make([]float32, len(positions))

	for i := range positions {
		policies[i] = make([]float32, PolicySize)
		policies[i][0] = 1.0
		values[i] = -0.99 // Return high losing score so parent gets +0.99
	}
	return policies, values, nil
}

func TestMCTSSearcher_EarlyExitMonitor(t *testing.T) {
	pos, _ := engine.ParseFEN("7k/4N1pp/8/2n3N1/2P3K1/8/8/8 w - - 0 1") // Mate in 1

	evaluator := &mockEarlyExitEvaluator{firstCall: true}
	searcher := NewMCTSSearcher(evaluator, 100000, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	_, err := searcher.Search(ctx, pos)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Search(): unexpected error: %v", err)
	}

	if duration >= 2*time.Second {
		t.Errorf("Search() duration = %v, want < %v", duration, 2*time.Second)
	}
	if duration < 100*time.Millisecond {
		t.Errorf("Search() duration = %v, want >= %v", duration, 100*time.Millisecond)
	}
}
