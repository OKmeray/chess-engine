package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain/engine"
)

// ErrWorkerPanic is returned when a background worker encounters a fatal panic.
var ErrWorkerPanic = errors.New("engine crashed critically in background worker")

// Evaluator defines the interface for Neural Network evaluation.
type Evaluator interface {
	// Evaluate takes a batch of positions and returns the policy and value vectors.
	EvaluateBatch(positions []engine.Position) (policies [][]float32, values []float32, err error)
}

// Node States
const (
	StateUnexpanded int32 = 0
	StateExpanding  int32 = 1
	StateExpanded   int32 = 2
)

// valueScale is used for Fixed-Point Arithmetic in MCTSNode.TotalValue.
const valueScale float32 = 1000000.0

// MCTSNode represents a perfectly packed node for the Memory Arena.
type MCTSNode struct {
	Visits        int32
	VirtualLoss   int32
	TotalValue    int64
	Prior         float32
	FirstChildIdx int32
	NumChildren   int16
	Move          engine.Move // The move that led to this node
	State         int32
}

// addValue atomically adds a scaled float32 to the TotalValue field.
func (n *MCTSNode) addValue(val float32) {
	scaled := int64(val * valueScale)
	atomic.AddInt64(&n.TotalValue, scaled)
}

// getValue atomically reads the float32 TotalValue.
func (n *MCTSNode) getValue() float32 {
	scaled := atomic.LoadInt64(&n.TotalValue)
	return float32(scaled) / valueScale
}

// EvalRequest packages a board position and a return channel for the NN batcher.
type EvalRequest struct {
	NodeIdx int32
	Pos     engine.Position
	RespCh  chan EvalResponse
}

// EvalResponse carries the evaluated policy and value back to a waiting worker.
type EvalResponse struct {
	Policy []float32
	Value  float32
	Err    error
}

// MCTSSearcher manages the state and workers for a concurrent Monte Carlo Tree Search.
type MCTSSearcher struct {
	evaluator  Evaluator
	nodes      []MCTSNode
	nodeCount  int32
	numWorkers int
	evalReqCh  chan EvalRequest
}

// NewMCTSSearcher initializes a concurrent Monte Carlo Tree Search engine.
func NewMCTSSearcher(evaluator Evaluator, maxNodes int, numWorkers int) *MCTSSearcher {
	return &MCTSSearcher{
		evaluator:  evaluator,
		nodes:      make([]MCTSNode, maxNodes),
		nodeCount:  0,
		evalReqCh:  make(chan EvalRequest, numWorkers*2),
		numWorkers: numWorkers,
	}
}

// Search is the main entry point to find the best move.
func (m *MCTSSearcher) Search(ctx context.Context, pos *engine.Position) (engine.Move, error) {
	// Initialize the memory arena with a fresh root node for the new search.
	atomic.StoreInt32(&m.nodeCount, 1)
	m.nodes[0] = MCTSNode{
		State: StateUnexpanded,
	}

	var fatalErr atomic.Pointer[error]
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	safeGo := func(f func()) {
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("%w: %v", ErrWorkerPanic, r)
					fatalErr.CompareAndSwap(nil, &err)
					cancel() // Stop all the workers
				}
			}()
			f()
		})
	}

	safeGo(func() {
		m.nnBatchWorker(ctx)
	})

	// Monitor search progress to abort early if a forced win is detected.
	safeGo(func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				root := &m.nodes[0]
				if atomic.LoadInt32(&root.State) != StateExpanded {
					continue
				}

				start := root.FirstChildIdx
				end := start + int32(root.NumChildren)

				for i := start; i < end; i++ {
					child := &m.nodes[i]
					visits := atomic.LoadInt32(&child.Visits)

					// If a move has > 100 visits and mathematically
					// converges to a 98%+ win rate, abort search early
					if visits > 100 {
						q := child.getValue() / float32(visits)
						if q >= 0.98 {
							cancel()
							return
						}
					}
				}
			}
		}
	})

	for i := 0; i < m.numWorkers; i++ {
		safeGo(func() {
			m.treeWorker(ctx, pos, &fatalErr, cancel)
		})
	}

	wg.Wait()

	if errPtr := fatalErr.Load(); errPtr != nil {
		return engine.Move(0), *errPtr
	}

	return m.getBestMove(0), nil
}

func (m *MCTSSearcher) treeWorker(ctx context.Context, rootPos *engine.Position, fatalErr *atomic.Pointer[error], cancel context.CancelFunc) {
	pos := &engine.Position{}
	moves := make([]engine.Move, 0, 200)
	history := make([]int32, 0, 256)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		currIdx := int32(0)
		history = history[:0]

		// Reset local position to root
		*pos = *rootPos

		// Select Leaf
		for {
			node := &m.nodes[currIdx]
			if atomic.LoadInt32(&node.State) != StateExpanded {
				break // Reached an unexpanded leaf
			}
			if node.NumChildren == 0 {
				break
			}

			// Apply virtual loss to discourage other workers
			m.applyVirtualLoss(currIdx)

			history = append(history, currIdx)
			currIdx = m.selectChild(currIdx)
			pos.MakeMove(m.nodes[currIdx].Move)
		}

		node := &m.nodes[currIdx]

		if atomic.LoadInt32(&node.State) == StateExpanded && node.NumChildren == 0 {
			// Apply virtual loss for the terminal leaf itself
			m.applyVirtualLoss(currIdx)
			history = append(history, currIdx)

			var val float32
			if pos.IsSquareAttacked(pos.GetKingSq(pos.SideToMove), 6-pos.SideToMove) {
				val = -1.0 // Checkmate
			} else {
				val = 0.0 // Stalemate
			}

			m.revertVirtualLoss(history)
			m.backpropagate(history, val)
			continue
		}

		if atomic.CompareAndSwapInt32(&node.State, StateUnexpanded, StateExpanding) {
			// Apply virtual loss for the leaf itself
			m.applyVirtualLoss(currIdx)
			history = append(history, currIdx)

			// NN Batch Request
			respCh := make(chan EvalResponse, 1)

			select {
			case m.evalReqCh <- EvalRequest{
				Pos:    *pos,
				RespCh: respCh,
			}:
			case <-ctx.Done():
				m.revertVirtualLoss(history)
				return
			}

			// Wait for NN response
			var resp EvalResponse
			select {
			case resp = <-respCh:
				if resp.Err != nil {
					m.revertVirtualLoss(history)
					fatalErr.CompareAndSwap(nil, &resp.Err)
					cancel()
					return
				}
			case <-ctx.Done():
				m.revertVirtualLoss(history)
				return
			}

			// Terminal Node check
			moves = moves[:0]
			moves = pos.GenerateMoves(moves)

			if len(moves) == 0 {
				atomic.StoreInt32(&node.State, StateExpanded)
				var val float32
				if pos.IsSquareAttacked(pos.GetKingSq(pos.SideToMove), 6-pos.SideToMove) {
					val = -1.0 // Checkmate
				} else {
					val = 0.0 // Stalemate
				}

				m.revertVirtualLoss(history)
				m.backpropagate(history, val)
				continue
			}

			m.expand(currIdx, &resp, pos.SideToMove, moves)

			m.revertVirtualLoss(history)
			m.backpropagate(history, resp.Value)

		} else {
			// Collision detected: another worker holds the expansion lock for this node.
			m.revertVirtualLoss(history)

			// Relinquish the CPU thread to allow the batcher to execute.
			runtime.Gosched()
			continue
		}
	}
}

func (m *MCTSSearcher) applyVirtualLoss(idx int32) {
	node := &m.nodes[idx]
	atomic.AddInt32(&node.VirtualLoss, 1)
}

func (m *MCTSSearcher) revertVirtualLoss(history []int32) {
	for _, idx := range history {
		node := &m.nodes[idx]
		atomic.AddInt32(&node.VirtualLoss, -1)
	}
}

func (m *MCTSSearcher) nnBatchWorker(ctx context.Context) {
	const maxBatchSize = 16
	batchPositions := make([]engine.Position, 0, maxBatchSize)
	batchChannels := make([]chan EvalResponse, 0, maxBatchSize)

	for {
		batchPositions = batchPositions[:0]
		batchChannels = batchChannels[:0]

		select {
		case <-ctx.Done():
			return
		case req := <-m.evalReqCh:
			batchPositions = append(batchPositions, req.Pos)
			batchChannels = append(batchChannels, req.RespCh)

			// Greedily grab more requests if available
		gatherLoop:
			for len(batchPositions) < maxBatchSize {
				select {
				case r := <-m.evalReqCh:
					batchPositions = append(batchPositions, r.Pos)
					batchChannels = append(batchChannels, r.RespCh)
				default:
					break gatherLoop
				}
			}

			// Evaluate Batch
			policies, values, err := m.evaluator.EvaluateBatch(batchPositions)
			if err != nil {
				// Propagate error to all waiting workers
				for i := 0; i < len(batchChannels); i++ {
					batchChannels[i] <- EvalResponse{Err: err}
				}
				continue
			}

			// Send responses back
			for i := 0; i < len(batchChannels); i++ {
				batchChannels[i] <- EvalResponse{
					Policy: policies[i],
					Value:  values[i],
				}
			}
		}
	}
}

func (m *MCTSSearcher) expand(nodeIdx int32, resp *EvalResponse, sideToMove engine.PieceColor, legalMoves []engine.Move) {
	numChildren := int32(len(legalMoves))
	if numChildren == 0 {
		atomic.StoreInt32(&m.nodes[nodeIdx].State, StateExpanded)
		return
	}

	// Allocate contiguous block of children in the arena
	firstChildIdx := atomic.AddInt32(&m.nodeCount, numChildren) - numChildren

	// if the arena is full, mark node as terminal
	if int(firstChildIdx+numChildren) > len(m.nodes) {
		atomic.AddInt32(&m.nodeCount, -numChildren) // Roll back allocation
		atomic.StoreInt32(&m.nodes[nodeIdx].State, StateExpanded)
		return
	}

	m.nodes[nodeIdx].FirstChildIdx = firstChildIdx
	m.nodes[nodeIdx].NumChildren = int16(numChildren)

	// Extract logits and apply Softmax over legal moves
	logits := make([]float32, len(legalMoves))
	maxLogit := float32(-math.MaxFloat32)

	for i, move := range legalMoves {
		moveIndex := EncodeMove(move, sideToMove == engine.Black)
		var logit float32 = -1000.0 // Default very low logit
		if moveIndex >= 0 && moveIndex < len(resp.Policy) {
			logit = resp.Policy[moveIndex]
		}

		logits[i] = logit

		if logit > maxLogit {
			maxLogit = logit
		}
	}

	sumExp := float32(0.0)
	for i := range logits {
		logits[i] = float32(math.Exp(float64(logits[i] - maxLogit)))
		sumExp += logits[i]
	}

	for i, move := range legalMoves {
		childIdx := firstChildIdx + int32(i)

		prior := logits[i] / sumExp

		m.nodes[childIdx] = MCTSNode{
			Prior: prior,
			Move:  move,
			State: StateUnexpanded,
		}
	}

	// Mark as fully expanded
	atomic.StoreInt32(&m.nodes[nodeIdx].State, StateExpanded)
}

// selectChild implements the AlphaZero PUCT (Predictor + Upper Confidence Bound) selection formula.
// The formula balances Exploitation vs Exploration:
// - Q (Action Value): Exploitation. The average evaluation score this node has yielded so far.
// - U (Confidence Bound): Exploration. A bonus based on the Neural Network's Prior
// probability, which decays as the node gets visited more often.
// The algorithm selects the child that maximizes (Q + U).
func (m *MCTSSearcher) selectChild(nodeIdx int32) int32 {
	node := &m.nodes[nodeIdx]
	parentVisits := atomic.LoadInt32(&node.Visits) + atomic.LoadInt32(&node.VirtualLoss)

	sqrtParentN := float32(math.Sqrt(float64(parentVisits)))

	var bestChild int32 = -1
	bestScore := float32(math.Inf(-1))

	start := node.FirstChildIdx
	end := start + int32(node.NumChildren)

	for i := start; i < end; i++ {
		child := &m.nodes[i]
		visits := atomic.LoadInt32(&child.Visits)
		vLoss := atomic.LoadInt32(&child.VirtualLoss)
		totalVisits := visits + vLoss

		var q float32
		if totalVisits > 0 {
			q = (child.getValue() - float32(vLoss)) / float32(totalVisits)
		}

		u := child.Prior * 1.5 * sqrtParentN / float32(1+totalVisits)
		score := q + u

		if score > bestScore {
			bestScore = score
			bestChild = i
		}
	}

	return bestChild
}

func (m *MCTSSearcher) backpropagate(history []int32, value float32) {
	// Root node is history[0]. Leaf is history[len-1].
	// Value should negate at each step upwards.
	currentVal := -value

	for i := len(history) - 1; i >= 0; i-- {
		idx := history[i]
		node := &m.nodes[idx]

		atomic.AddInt32(&node.Visits, 1)
		node.addValue(currentVal)

		currentVal = -currentVal
	}
}

func (m *MCTSSearcher) getBestMove(rootIdx int32) engine.Move {
	node := &m.nodes[rootIdx]

	start := node.FirstChildIdx
	end := start + int32(node.NumChildren)

	var bestMove engine.Move
	maxVisits := int32(-1)

	for i := start; i < end; i++ {
		child := &m.nodes[i]
		visits := atomic.LoadInt32(&child.Visits)

		if visits > maxVisits {
			maxVisits = visits
			bestMove = child.Move
		}
	}

	return bestMove
}
