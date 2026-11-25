package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OKmeray/chess-engine/internal/domain"
	"github.com/OKmeray/chess-engine/internal/transport/rest"
	"github.com/OKmeray/chess-engine/internal/transport/ws"
	"github.com/OKmeray/chess-engine/internal/usecase"
	"go.uber.org/zap"
)

type mockModelService struct{}

func (m *mockModelService) List(_ context.Context, _, _ int) ([]domain.NNModel, int, error) {
	return []domain.NNModel{}, 0, nil
}

type mockGameCreator struct{}

func (m *mockGameCreator) CreateGame(_ context.Context, _ usecase.EngineType, _ int) (usecase.MoveFinder, error) {
	return nil, nil
}

func TestNewRouter(t *testing.T) {
	logger := zap.NewNop().Sugar()

	t.Run("MountsBothHandlers", func(t *testing.T) {
		modelHandler := rest.NewModelHandler(&mockModelService{}, logger)
		engineHandler := ws.NewEngineHandler(&mockGameCreator{}, logger)

		mux := NewRouter(RouterConfig{
			ModelHandler:  modelHandler,
			EngineHandler: engineHandler,
		})

		// Verify REST endpoint is registered
		restReq := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
		restRR := httptest.NewRecorder()
		mux.ServeHTTP(restRR, restReq)
		if restRR.Code != http.StatusOK {
			t.Errorf("REST route /api/models: got %d, want %d", restRR.Code, http.StatusOK)
		}

		// Verify WebSocket endpoint is registered
		wsReq := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ws/game", nil)
		wsRR := httptest.NewRecorder()
		mux.ServeHTTP(wsRR, wsReq)
		if wsRR.Code != http.StatusUpgradeRequired {
			t.Errorf("WS route /ws/game: got %d, want %d", wsRR.Code, http.StatusUpgradeRequired)
		}
	})

	t.Run("NilHandlersAllowed", func(t *testing.T) {
		mux := NewRouter(RouterConfig{})
		if mux == nil {
			t.Fatalf("NewRouter() = nil, want non-nil *http.ServeMux")
		}

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}
