package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/OKmeray/chess-engine/internal/domain"
	"go.uber.org/zap"
)

type mockModelService struct {
	models        []domain.NNModel
	total         int
	err           error
	recordedLimit int
	recordedSkip  int
	calls         int
}

func (m *mockModelService) List(_ context.Context, limit, offset int) ([]domain.NNModel, int, error) {
	m.calls++
	m.recordedLimit = limit
	m.recordedSkip = offset
	return m.models, m.total, m.err
}

func TestNewModelHandler(t *testing.T) {
	logger := zap.NewNop().Sugar()
	svc := &mockModelService{}

	handler := NewModelHandler(svc, logger)
	if handler == nil {
		t.Fatalf("NewModelHandler() = nil, want non-nil ModelHandler")
	}
	if handler.svc != svc {
		t.Errorf("handler.svc = %v, want %v", handler.svc, svc)
	}
	if handler.logger != logger {
		t.Errorf("handler.logger = %v, want %v", handler.logger, logger)
	}
}

func TestModelHandler_RegisterRoutes(t *testing.T) {
	logger := zap.NewNop().Sugar()
	svc := &mockModelService{}
	handler := NewModelHandler(svc, logger)
	mux := http.NewServeMux()

	handler.RegisterRoutes(mux)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("RegisterRoutes: rr.Code = %d, want %d", rr.Code, http.StatusOK)
	}

	postReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/models", nil)
	postRR := httptest.NewRecorder()
	mux.ServeHTTP(postRR, postReq)

	if postRR.Code != http.StatusMethodNotAllowed {
		t.Errorf("RegisterRoutes POST: postRR.Code = %d, want %d", postRR.Code, http.StatusMethodNotAllowed)
	}
	if allow := postRR.Header().Get("Allow"); allow != "GET, HEAD" {
		t.Errorf("RegisterRoutes POST Allow: got %q, want %q", allow, "GET, HEAD")
	}
}

func TestModelHandler_HandleListModels_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	sampleTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	expectedModels := []domain.NNModel{
		{
			ID:        1,
			Name:      "v1-cnn",
			Path:      "/models/v1.onnx",
			IsActive:  true,
			CreatedAt: sampleTime,
			Details: domain.NNModelDetails{
				Name:         "v1-cnn",
				Architecture: "CNN",
			},
		},
		{
			ID:        2,
			Name:      "v2-transformer",
			Path:      "/models/v2.onnx",
			IsActive:  false,
			CreatedAt: sampleTime,
			Details: domain.NNModelDetails{
				Name:         "v2-transformer",
				Architecture: "Transformer",
			},
		},
	}

	svc := &mockModelService{
		models: expectedModels,
		total:  len(expectedModels),
	}
	handler := NewModelHandler(svc, logger)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
	rr := httptest.NewRecorder()

	handler.HandleListModels(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var resp PaginatedResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if resp.Total != len(expectedModels) {
		t.Errorf("resp.Total = %d, want %d", resp.Total, len(expectedModels))
	}
	if resp.Limit != defaultLimit {
		t.Errorf("resp.Limit = %d, want %d", resp.Limit, defaultLimit)
	}
	if resp.Skip != 0 {
		t.Errorf("resp.Skip = %d, want %d", resp.Skip, 0)
	}
	if len(resp.Data) != len(expectedModels) {
		t.Errorf("Len: got %d, want %d", len(resp.Data), len(expectedModels))
	}
	if !reflect.DeepEqual(resp.Data, expectedModels) {
		t.Errorf("resp.Data:\ngot:  %#v\nwant: %#v", resp.Data, expectedModels)
	}
}

func TestModelHandler_HandleListModels_PaginationParams(t *testing.T) {
	logger := zap.NewNop().Sugar()

	tests := []struct {
		name      string
		rawQuery  string
		wantLimit int
		wantSkip  int
	}{
		{
			name:      "Default params when omitted",
			rawQuery:  "",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Custom valid limit and skip",
			rawQuery:  "?limit=25&skip=50",
			wantLimit: 25,
			wantSkip:  50,
		},
		{
			name:      "Exact maximum limit is accepted",
			rawQuery:  fmt.Sprintf("?limit=%d", maxLimit),
			wantLimit: maxLimit,
			wantSkip:  0,
		},
		{
			name:      "Limit exceeding maximum is clamped",
			rawQuery:  fmt.Sprintf("?limit=%d", maxLimit+1),
			wantLimit: maxLimit,
			wantSkip:  0,
		},
		{
			name:      "Lowest positive limit is accepted",
			rawQuery:  "?limit=1",
			wantLimit: 1,
			wantSkip:  0,
		},
		{
			name:      "Zero limit falls back to default",
			rawQuery:  "?limit=0",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Negative limit falls back to default",
			rawQuery:  "?limit=-1",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Non-numeric limit falls back to default",
			rawQuery:  "?limit=invalid",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Zero skip is valid",
			rawQuery:  "?skip=0",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Negative skip falls back to default",
			rawQuery:  "?skip=-1",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
		{
			name:      "Non-numeric skip falls back to default",
			rawQuery:  "?skip=abc",
			wantLimit: defaultLimit,
			wantSkip:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockModelService{
				models: []domain.NNModel{},
				total:  0,
			}
			handler := NewModelHandler(svc, logger)

			targetURL := "/api/models" + tt.rawQuery
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, targetURL, nil)
			rr := httptest.NewRecorder()

			handler.HandleListModels(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusOK)
			}
			if svc.recordedLimit != tt.wantLimit {
				t.Errorf("svc.recordedLimit = %d, want %d", svc.recordedLimit, tt.wantLimit)
			}
			if svc.recordedSkip != tt.wantSkip {
				t.Errorf("svc.recordedSkip = %d, want %d", svc.recordedSkip, tt.wantSkip)
			}
		})
	}
}

func TestModelHandler_HandleListModels_ServiceError(t *testing.T) {
	logger := zap.NewNop().Sugar()
	svc := &mockModelService{
		err: errors.New("database connection failure"),
	}
	handler := NewModelHandler(svc, logger)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
	rr := httptest.NewRecorder()

	handler.HandleListModels(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if svc.calls != 1 {
		t.Errorf("svc.calls = %d, want %d", svc.calls, 1)
	}
}

func TestModelHandler_HandleListModels_NilModelsReturnsEmptySlice(t *testing.T) {
	logger := zap.NewNop().Sugar()
	svc := &mockModelService{
		models: nil,
		total:  0,
	}
	handler := NewModelHandler(svc, logger)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
	rr := httptest.NewRecorder()

	handler.HandleListModels(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusOK)
	}

	var resp PaginatedResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if resp.Data == nil {
		t.Errorf("resp.Data = nil, want empty slice []")
	}
	if len(resp.Data) != 0 {
		t.Errorf("Len: got %d, want %d", len(resp.Data), 0)
	}
}

func TestModelHandler_HandleListModels_NilService(t *testing.T) {
	logger := zap.NewNop().Sugar()
	handler := NewModelHandler(nil, logger)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/models", nil)
	rr := httptest.NewRecorder()

	handler.HandleListModels(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("rr.Code = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}
}
