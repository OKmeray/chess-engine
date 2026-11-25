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

type mockMoveFinder struct {
	mu                sync.Mutex
	recordedFEN       string
	recordedRemaining time.Duration
	recordedIncrement time.Duration
}

func (m *mockMoveFinder) FindBestMove(_ context.Context, fen string, remainingTime, increment time.Duration) (string, error) {
	m.mu.Lock()
	m.recordedFEN = fen
	m.recordedRemaining = remainingTime
	m.recordedIncrement = increment
	m.mu.Unlock()
	return "e2e4", nil
}

type mockGameCreator struct {
	finder       usecase.MoveFinder
	err          error
	recordedType usecase.EngineType
	recordedID   int
}

func (m *mockGameCreator) CreateGame(_ context.Context, engineType usecase.EngineType, botID int) (usecase.MoveFinder, error) {
	m.recordedType = engineType
	m.recordedID = botID
	return m.finder, m.err
}

func TestEngineHandler_ServeWS_OriginPolicy(t *testing.T) {
	logger := zap.NewNop().Sugar()

	tests := []struct {
		name           string
		originPatterns []string
		origin         string
		wantAllowed    bool
	}{
		{
			name:        "Default allows localhost",
			origin:      "http://localhost:3000",
			wantAllowed: true,
		},
		{
			name:        "Default allows 127.0.0.1",
			origin:      "http://127.0.0.1:5173",
			wantAllowed: true,
		},
		{
			name:        "Default rejects untrusted origin",
			origin:      "https://evil.com",
			wantAllowed: false,
		},
		{
			name:           "Custom allows configured origin",
			originPatterns: []string{"custom-origin.com"},
			origin:         "https://custom-origin.com",
			wantAllowed:    true,
		},
		{
			name:           "Custom rejects localhost when overridden",
			originPatterns: []string{"custom-origin.com"},
			origin:         "http://localhost:3000",
			wantAllowed:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewEngineHandler(&mockGameCreator{}, logger, tt.originPatterns...)
			srv := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
			t.Cleanup(srv.Close)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			t.Cleanup(cancel)

			wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
			conn, resp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
				HTTPHeader: http.Header{"Origin": {tt.origin}},
			})
			if conn != nil {
				t.Cleanup(func() {
					_ = conn.Close(websocket.StatusNormalClosure, "")
				})
			}

			if tt.wantAllowed {
				if err != nil {
					t.Fatalf("Dial(origin=%q): unexpected error: %v", tt.origin, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Dial(origin=%q) error = nil, wantErr %v", tt.origin, true)
			}
			if resp == nil || resp.StatusCode != http.StatusForbidden {
				gotCode := 0
				if resp != nil {
					gotCode = resp.StatusCode
				}
				t.Errorf("Dial(origin=%q) status = %d, want %d", tt.origin, gotCode, http.StatusForbidden)
			}
		})
	}
}

func TestEngineHandler_RegisterRoutes(t *testing.T) {
	logger := zap.NewNop().Sugar()
	handler := NewEngineHandler(&mockGameCreator{}, logger)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ws/game", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUpgradeRequired {
		t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusUpgradeRequired)
	}
}

func TestEngineHandler_ServeWS_NotConfigured(t *testing.T) {
	logger := zap.NewNop().Sugar()
	var nilHandler *EngineHandler

	tests := []struct {
		name    string
		handler *EngineHandler
	}{
		{
			name:    "Nil handler receiver",
			handler: nilHandler,
		},
		{
			name:    "Nil game creator",
			handler: NewEngineHandler(nil, logger),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ws/game", nil)
			rr := httptest.NewRecorder()

			tt.handler.ServeWS(rr, req)

			if rr.Code != http.StatusInternalServerError {
				t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusInternalServerError)
			}
		})
	}
}

func TestEngineHandler_ServeWS_ValidConnection(t *testing.T) {
	logger := zap.NewNop().Sugar()
	creator := &mockGameCreator{
		finder: &mockMoveFinder{},
	}

	handler := NewEngineHandler(creator, logger)

	srv := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	conn := dialTestClient(t, ctx, wsURL)

	initMsg := ClientMessage{
		Type:       MessageTypeInitGame,
		EngineType: string(usecase.EngineTypeMCTS),
		ModelID:    1,
	}
	if err := wsjson.Write(ctx, conn, initMsg); err != nil {
		t.Fatalf("failed to write init_game message: %v", err)
	}

	var resp ServerMessage
	if err := wsjson.Read(ctx, conn, &resp); err != nil {
		t.Fatalf("failed to read ready message: %v", err)
	}

	if resp.Type != MessageTypeReady {
		t.Errorf("resp.Type = %v, want %v", resp.Type, MessageTypeReady)
	}
	if resp.EngineType != string(usecase.EngineTypeMCTS) {
		t.Errorf("resp.EngineType = %q, want %q", resp.EngineType, usecase.EngineTypeMCTS)
	}
	if resp.ModelID != 1 {
		t.Errorf("resp.ModelID = %d, want %d", resp.ModelID, 1)
	}
	if creator.recordedType != usecase.EngineTypeMCTS {
		t.Errorf("creator.recordedType = %q, want %q", creator.recordedType, usecase.EngineTypeMCTS)
	}
	if creator.recordedID != 1 {
		t.Errorf("creator.recordedID = %d, want %d", creator.recordedID, 1)
	}
}
