package ws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/usecase"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"go.uber.org/zap"
)

// dialTestClient connects to a test WebSocket server and registers automatic cleanup.
func dialTestClient(t *testing.T, ctx context.Context, wsURL string) *websocket.Conn {
	t.Helper()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	})
	return conn
}

// setupTestClient boots a test HTTP server hosting a Client and returns an open connection.
func setupTestClient(t *testing.T, ctx context.Context, creator GameCreator, finder MoveFinder) *websocket.Conn {
	t.Helper()
	logger := zap.NewNop().Sugar()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Errorf("failed to accept websocket: %v", err)
			return
		}
		client := NewClient(conn, creator, logger)
		client.finder = finder
		client.Listen(r.Context())
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	return dialTestClient(t, ctx, wsURL)
}

type errorMoveFinder struct{}

func (m *errorMoveFinder) FindBestMove(_ context.Context, _ string, _, _ time.Duration) (string, error) {
	return "", usecase.ErrSearchFailed
}

func TestClient_MessageHandling(t *testing.T) {
	tests := []struct {
		name           string
		creator        GameCreator
		finder         MoveFinder
		sendMsg        ClientMessage
		wantType       MessageType
		wantBestMove   string
		wantFEN        string
		wantRemaining  time.Duration
		wantIncrement  time.Duration
		wantEngineType string
		wantModelID    int
		wantErrMsg     string
	}{
		{
			name:   "Calculate move returns best move",
			finder: &mockMoveFinder{},
			sendMsg: ClientMessage{
				Type:          MessageTypeCalculateMove,
				FEN:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
				RemainingTime: 60000,
				Increment:     3000,
			},
			wantType:      MessageTypeBestMove,
			wantBestMove:  "e2e4",
			wantFEN:       "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			wantRemaining: 60 * time.Second,
			wantIncrement: 3 * time.Second,
		},
		{
			name:       "Calculate move error returns error message",
			finder:     &errorMoveFinder{},
			sendMsg:    ClientMessage{Type: MessageTypeCalculateMove},
			wantType:   MessageTypeError,
			wantErrMsg: ErrCalculationFailed.Error(),
		},
		{
			name:       "Calculate move when finder is nil returns error",
			sendMsg:    ClientMessage{Type: MessageTypeCalculateMove},
			wantType:   MessageTypeError,
			wantErrMsg: ErrEngineNotInitialized.Error(),
		},
		{
			name:           "Init game switches engine successfully",
			creator:        &mockGameCreator{finder: &mockMoveFinder{}},
			sendMsg:        ClientMessage{Type: MessageTypeInitGame, EngineType: string(usecase.EngineTypeMCTS), ModelID: 2},
			wantType:       MessageTypeReady,
			wantEngineType: string(usecase.EngineTypeMCTS),
			wantModelID:    2,
		},
		{
			name:           "Init game with default engine type",
			creator:        &mockGameCreator{finder: &mockMoveFinder{}},
			sendMsg:        ClientMessage{Type: MessageTypeInitGame},
			wantType:       MessageTypeReady,
			wantEngineType: string(usecase.EngineTypeMCTS),
		},
		{
			name:       "Init game with unsupported engine returns error",
			creator:    &mockGameCreator{err: usecase.ErrUnsupportedEngineType},
			sendMsg:    ClientMessage{Type: MessageTypeInitGame, EngineType: "unsupported"},
			wantType:   MessageTypeError,
			wantErrMsg: ErrInitGameFailed.Error(),
		},
		{
			name:       "Init game when creator is nil returns error",
			sendMsg:    ClientMessage{Type: MessageTypeInitGame},
			wantType:   MessageTypeError,
			wantErrMsg: ErrCreatorNotConfigured.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()

			conn := setupTestClient(t, ctx, tt.creator, tt.finder)

			if err := wsjson.Write(ctx, conn, tt.sendMsg); err != nil {
				t.Fatalf("failed to write message: %v", err)
			}

			var resp ServerMessage
			if err := wsjson.Read(ctx, conn, &resp); err != nil {
				t.Fatalf("failed to read response: %v", err)
			}

			if resp.Type != tt.wantType {
				t.Errorf("resp.Type = %v, want %v", resp.Type, tt.wantType)
			}
			if tt.wantBestMove != "" && resp.BestMove != tt.wantBestMove {
				t.Errorf("resp.BestMove = %q, want %q", resp.BestMove, tt.wantBestMove)
			}
			if tt.wantEngineType != "" && resp.EngineType != tt.wantEngineType {
				t.Errorf("resp.EngineType = %q, want %q", resp.EngineType, tt.wantEngineType)
			}
			if tt.wantType == MessageTypeReady && resp.ModelID != tt.wantModelID {
				t.Errorf("resp.ModelID = %d, want %d", resp.ModelID, tt.wantModelID)
			}
			if tt.wantErrMsg != "" && resp.Error != tt.wantErrMsg {
				t.Errorf("resp.Error = %q, want %q", resp.Error, tt.wantErrMsg)
			}
			if mockFinder, ok := tt.finder.(*mockMoveFinder); ok && tt.wantType == MessageTypeBestMove {
				mockFinder.mu.Lock()
				gotFEN := mockFinder.recordedFEN
				gotRemaining := mockFinder.recordedRemaining
				gotIncrement := mockFinder.recordedIncrement
				mockFinder.mu.Unlock()

				if gotFEN != tt.wantFEN {
					t.Errorf("finder fen = %q, want %q", gotFEN, tt.wantFEN)
				}
				if gotRemaining != tt.wantRemaining {
					t.Errorf("finder remainingTime = %v, want %v", gotRemaining, tt.wantRemaining)
				}
				if gotIncrement != tt.wantIncrement {
					t.Errorf("finder increment = %v, want %v", gotIncrement, tt.wantIncrement)
				}
			}
		})
	}
}

func TestClient_InitGame_ThenCalculateMove(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	creator := &mockGameCreator{finder: &mockMoveFinder{}}
	conn := setupTestClient(t, ctx, creator, nil)

	initMsg := ClientMessage{
		Type:       MessageTypeInitGame,
		EngineType: string(usecase.EngineTypeMinimax),
	}
	if err := wsjson.Write(ctx, conn, initMsg); err != nil {
		t.Fatalf("failed to write init message: %v", err)
	}

	var initResp ServerMessage
	if err := wsjson.Read(ctx, conn, &initResp); err != nil {
		t.Fatalf("failed to read init response: %v", err)
	}
	if initResp.Type != MessageTypeReady {
		t.Errorf("initResp.Type = %v, want %v", initResp.Type, MessageTypeReady)
	}

	calcMsg := ClientMessage{Type: MessageTypeCalculateMove}
	if err := wsjson.Write(ctx, conn, calcMsg); err != nil {
		t.Fatalf("failed to write calculate message: %v", err)
	}

	var moveResp ServerMessage
	if err := wsjson.Read(ctx, conn, &moveResp); err != nil {
		t.Fatalf("failed to read move response: %v", err)
	}
	if moveResp.Type != MessageTypeBestMove {
		t.Errorf("moveResp.Type = %v, want %v", moveResp.Type, MessageTypeBestMove)
	}
}

func TestClient_UnknownMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	conn := setupTestClient(t, ctx, nil, &mockMoveFinder{})

	req := ClientMessage{Type: "unknown_type"}
	if err := wsjson.Write(ctx, conn, req); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	// Verify the client didn't crash by sending a valid message right after
	validReq := ClientMessage{Type: MessageTypeCalculateMove}
	if err := wsjson.Write(ctx, conn, validReq); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	var resp ServerMessage
	if err := wsjson.Read(ctx, conn, &resp); err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if resp.Type != MessageTypeBestMove {
		t.Errorf("resp.Type = %v, want %v", resp.Type, MessageTypeBestMove)
	}
}

func TestClient_MalformedJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	conn := setupTestClient(t, ctx, nil, &mockMoveFinder{})

	if err := conn.Write(ctx, websocket.MessageText, []byte("{bad-json")); err != nil {
		t.Fatalf("failed to write malformed message: %v", err)
	}

	var errResp ServerMessage
	if err := wsjson.Read(ctx, conn, &errResp); err != nil {
		t.Fatalf("failed to read error response: %v", err)
	}
	if errResp.Type != MessageTypeError {
		t.Errorf("errResp.Type = %v, want %v", errResp.Type, MessageTypeError)
	}
	if errResp.Error != ErrInvalidMessage.Error() {
		t.Errorf("errResp.Error = %q, want %q", errResp.Error, ErrInvalidMessage.Error())
	}

	// Verify the connection stayed open and processes subsequent valid requests
	if err := wsjson.Write(ctx, conn, ClientMessage{Type: MessageTypeCalculateMove}); err != nil {
		t.Fatalf("failed to write follow-up message: %v", err)
	}

	var moveResp ServerMessage
	if err := wsjson.Read(ctx, conn, &moveResp); err != nil {
		t.Fatalf("failed to read follow-up response: %v", err)
	}
	if moveResp.Type != MessageTypeBestMove {
		t.Errorf("moveResp.Type = %v, want %v", moveResp.Type, MessageTypeBestMove)
	}
}

type blockingMoveFinder struct {
	once    sync.Once
	started chan struct{}
	wait    chan struct{}
}

func (m *blockingMoveFinder) FindBestMove(ctx context.Context, _ string, _, _ time.Duration) (string, error) {
	if m.started != nil {
		m.once.Do(func() {
			close(m.started)
		})
	}
	select {
	case <-m.wait:
		return "e2e4", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func TestClient_EngineBusy(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	blocker := &blockingMoveFinder{
		started: make(chan struct{}),
		wait:    make(chan struct{}),
	}
	defer close(blocker.wait)

	conn := setupTestClient(t, ctx, nil, blocker)

	// Send first request, it goes to worker and blocks
	req1 := ClientMessage{Type: MessageTypeCalculateMove}
	if err := wsjson.Write(ctx, conn, req1); err != nil {
		t.Fatalf("failed to write first message: %v", err)
	}

	// Wait until worker picks up req1 deterministically without time.Sleep
	select {
	case <-blocker.started:
	case <-ctx.Done():
		t.Fatalf("timed out waiting for worker to start: %v", ctx.Err())
	}

	// Send second request, this goes into the size-1 buffer
	req2 := ClientMessage{Type: MessageTypeCalculateMove}
	if err := wsjson.Write(ctx, conn, req2); err != nil {
		t.Fatalf("failed to write second message: %v", err)
	}

	// Send third request, this should be dropped and error returned immediately
	req3 := ClientMessage{Type: MessageTypeCalculateMove}
	if err := wsjson.Write(ctx, conn, req3); err != nil {
		t.Fatalf("failed to write third message: %v", err)
	}

	var resp ServerMessage
	if err := wsjson.Read(ctx, conn, &resp); err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if resp.Type != MessageTypeError {
		t.Errorf("resp.Type = %v, want %v", resp.Type, MessageTypeError)
	}
	if resp.Error != ErrEngineBusy.Error() {
		t.Errorf("resp.Error = %q, want %q", resp.Error, ErrEngineBusy.Error())
	}
}
