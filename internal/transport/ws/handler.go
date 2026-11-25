package ws

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"go.uber.org/zap"

	"github.com/OKmeray/chess-engine/internal/usecase"
)

// GameCreator instantiates a MoveFinder per connection.
type GameCreator interface {
	CreateGame(ctx context.Context, engineType usecase.EngineType, botID int) (usecase.MoveFinder, error)
}

// EngineHandler handles WebSocket upgrades and client routing.
type EngineHandler struct {
	creator        GameCreator
	logger         *zap.SugaredLogger
	originPatterns []string
}

// NewEngineHandler initializes the EngineHandler.
func NewEngineHandler(creator GameCreator, logger *zap.SugaredLogger, originPatterns ...string) *EngineHandler {
	patterns := originPatterns
	if len(patterns) == 0 {
		patterns = []string{"localhost:*", "127.0.0.1:*"}
	}
	return &EngineHandler{
		creator:        creator,
		logger:         logger,
		originPatterns: patterns,
	}
}

// RegisterRoutes registers the WebSocket endpoints with the provided HTTP router.
func (h *EngineHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/game", h.ServeWS)
}

// ServeWS upgrades the HTTP connection and begins client processing.
func (h *EngineHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.creator == nil {
		if h != nil && h.logger != nil {
			h.logger.Error("game creator is not configured")
		}
		http.Error(w, ErrCreatorNotConfigured.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infow("incoming websocket connection", "remote_addr", r.RemoteAddr)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: h.originPatterns,
	})
	if err != nil {
		h.logger.Errorw("failed to upgrade connection", "error", err)
		return
	}

	client := NewClient(conn, h.creator, h.logger)
	client.Listen(r.Context())
}
