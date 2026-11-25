package transport

import (
	"net/http"

	"github.com/OKmeray/chess-engine/internal/transport/rest"
	"github.com/OKmeray/chess-engine/internal/transport/ws"
)

// RouterConfig bundles the transport handlers required to mount all HTTP endpoints.
type RouterConfig struct {
	ModelHandler  *rest.ModelHandler
	EngineHandler *ws.EngineHandler
}

// NewRouter constructs a unified HTTP multiplexer with all REST and WebSocket routes.
func NewRouter(cfg RouterConfig) *http.ServeMux {
	mux := http.NewServeMux()

	if cfg.ModelHandler != nil {
		cfg.ModelHandler.RegisterRoutes(mux)
	}
	if cfg.EngineHandler != nil {
		cfg.EngineHandler.RegisterRoutes(mux)
	}

	return mux
}
