package httpserver

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"mbsmf/internal/config"
	"mbsmf/internal/metrics"
	"mbsmf/internal/models"
	"mbsmf/internal/service"
	memstore "mbsmf/internal/store/memory"
)

type noNotify struct{}

func (n *noNotify) DispatchSubscriptionEvent(_ models.Subscription, _ string, _ any) {}

func newTestRouter() http.Handler {
	store := memstore.NewMemoryStore()
	m := metrics.NewRegistry()
	logger := log.New(os.Stdout, "test ", 0)
	svc := service.NewMBService(store, &noNotify{}, logger, m, config.LoadFromEnv())
	return NewRouter(svc, logger, m, config.LoadFromEnv())
}

func TestHealthz(t *testing.T) {
	h := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestSessionLifecycle(t *testing.T) {
	h := newTestRouter()
	// create session
	body, _ := json.Marshal(map[string]any{"tmgi": "10", "serviceArea": []string{"A"}})
	req := httptest.NewRequest(http.MethodPost, "/nmbsf/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
	// list
	req2 := httptest.NewRequest(http.MethodGet, "/nmbsf/v1/sessions", nil)
	resp2 := httptest.NewRecorder()
	h.ServeHTTP(resp2, req2)
	if resp2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.Code)
	}
}

