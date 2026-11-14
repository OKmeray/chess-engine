package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/OKmeray/chess-engine/internal/domain"
	"go.uber.org/zap"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

// ModelService defines the business logic interface for models.
type ModelService interface {
	List(ctx context.Context, limit, offset int) ([]domain.NNModel, int, error)
}

// ModelHandler handles HTTP REST requests for models.
type ModelHandler struct {
	svc    ModelService
	logger *zap.SugaredLogger
}

// NewModelHandler creates a new REST handler for models.
func NewModelHandler(svc ModelService, logger *zap.SugaredLogger) *ModelHandler {
	return &ModelHandler{
		svc:    svc,
		logger: logger,
	}
}

// RegisterRoutes registers the model endpoints with the provided HTTP router.
func (h *ModelHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/models", h.HandleListModels)
}

// PaginatedResponse standardizes the JSON structure for paginated results.
type PaginatedResponse struct {
	Data  []domain.NNModel `json:"data"`
	Total int              `json:"total"`
	Limit int              `json:"limit"`
	Skip  int              `json:"skip"`
}

// HandleListModels processes the GET request for listing paginated models.
func (h *ModelHandler) HandleListModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	if h == nil || h.svc == nil {
		http.Error(w, ErrServiceUnavailable.Error(), http.StatusServiceUnavailable)
		return
	}

	limit := defaultLimit
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = min(parsed, maxLimit)
		}
	}

	skip := 0
	if s := r.URL.Query().Get("skip"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			skip = parsed
		}
	}

	models, total, err := h.svc.List(r.Context(), limit, skip)
	if err != nil {
		h.logger.Errorw("failed to list models", "error", err)
		http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	if models == nil {
		models = []domain.NNModel{}
	}

	resp := PaginatedResponse{
		Data:  models,
		Total: total,
		Limit: limit,
		Skip:  skip,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Errorw("failed to encode response", "error", err)
	}
}
