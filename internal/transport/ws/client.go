package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/OKmeray/chess-engine/internal/usecase"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"go.uber.org/zap"
)

// MoveFinder finds the best move for a given position.
type MoveFinder interface {
	FindBestMove(ctx context.Context, fen string, remainingTime, increment time.Duration) (string, error)
}

// Client wraps a websocket connection for a single user.
type Client struct {
	conn    *websocket.Conn
	creator GameCreator
	finder  MoveFinder
	logger  *zap.SugaredLogger
	mu      sync.Mutex // Protects concurrent writes to conn and finder updates
	workQ   chan ClientMessage
}

// NewClient initializes a WebSocket client.
func NewClient(conn *websocket.Conn, creator GameCreator, logger *zap.SugaredLogger) *Client {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	return &Client{
		conn:    conn,
		creator: creator,
		logger:  logger,
		workQ:   make(chan ClientMessage, 1),
	}
}

// Listen blocks and routes incoming messages until the connection drops.
func (c *Client) Listen(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer wg.Wait()
	defer cancel()
	defer c.conn.Close(websocket.StatusNormalClosure, "client disconnected")

	wg.Go(func() {
		c.worker(ctx)
	})

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			c.logger.Infow("websocket read error or disconnect", "error", err)
			break
		}

		var msg ClientMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Warnw("invalid websocket message format", "error", err)
			c.sendError(ctx, ErrInvalidMessage.Error())
			continue
		}

		switch msg.Type {
		case MessageTypeCalculateMove:
			select {
			case c.workQ <- msg:
			default:
				c.logger.Warn("client requested calculate_move while engine is already busy")
				c.sendError(ctx, ErrEngineBusy.Error())
			}
		case MessageTypeInitGame:
			c.handleInitGame(ctx, msg)
		default:
			c.logger.Warnw("unknown message type received", "type", msg.Type)
		}
	}
}

func (c *Client) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.workQ:
			c.handleCalculateMove(ctx, msg)
		}
	}
}

func (c *Client) handleInitGame(ctx context.Context, msg ClientMessage) {
	if c.creator == nil {
		c.logger.Error("game creator is not configured")
		c.sendError(ctx, ErrCreatorNotConfigured.Error())
		return
	}

	engineTypeStr := msg.EngineType
	if engineTypeStr == "" {
		engineTypeStr = string(usecase.EngineTypeMCTS)
	}
	engineType := usecase.EngineType(engineTypeStr)

	newFinder, err := c.creator.CreateGame(ctx, engineType, msg.ModelID)
	if err != nil {
		c.logger.Errorw("failed to initialize game", "error", err)
		c.sendError(ctx, ErrInitGameFailed.Error())
		return
	}

	c.mu.Lock()
	c.finder = newFinder
	c.mu.Unlock()

	c.logger.Infow("game initialized", "engine_type", engineType, "model_id", msg.ModelID)
	_ = c.writeJSON(ctx, ServerMessage{
		Type:       MessageTypeReady,
		EngineType: string(engineType),
		ModelID:    msg.ModelID,
	})
}

func (c *Client) handleCalculateMove(ctx context.Context, msg ClientMessage) {
	c.mu.Lock()
	finder := c.finder
	c.mu.Unlock()

	if finder == nil {
		c.logger.Warn("no move finder configured")
		c.sendError(ctx, ErrEngineNotInitialized.Error())
		return
	}

	remainingTime := time.Duration(msg.RemainingTime) * time.Millisecond
	increment := time.Duration(msg.Increment) * time.Millisecond
	bestMove, err := finder.FindBestMove(ctx, msg.FEN, remainingTime, increment)
	if err != nil {
		c.logger.Errorw("failed to calculate move", "error", err)
		c.sendError(ctx, ErrCalculationFailed.Error())
		return
	}

	_ = c.writeJSON(ctx, ServerMessage{
		Type:     MessageTypeBestMove,
		BestMove: bestMove,
	})
}

func (c *Client) writeJSON(ctx context.Context, v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := wsjson.Write(ctx, c.conn, v); err != nil {
		c.logger.Errorw("failed to write websocket message", "error", err)
		return err
	}
	return nil
}

func (c *Client) sendError(ctx context.Context, errMsg string) {
	_ = c.writeJSON(ctx, ServerMessage{
		Type:  MessageTypeError,
		Error: errMsg,
	})
}
