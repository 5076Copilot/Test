package sbi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/mb-smf/internal/config"
	"example.com/mb-smf/internal/domain"
	"example.com/mb-smf/internal/pfcp"
	"example.com/mb-smf/internal/repo"
	"example.com/mb-smf/internal/service"
)

type Server struct {
	cfg    config.Config
	http   *http.Server
	svc    *service.Service
}

func NewServer(cfg config.Config) *Server {
	repository := repo.NewMemorySessionRepo()
	pfcpClient := pfcp.NewClient(cfg.PFCPAddr)
	svc := service.New(repository, pfcpClient)

	mux := http.NewServeMux()
	api := &apiHandler{svc: svc}
	mux.HandleFunc("/mbs/sessions", api.handleSessions)
	mux.HandleFunc("/mbs/sessions/", api.handleSessionByID)

	return &Server{cfg: cfg, http: &http.Server{Addr: cfg.Address, Handler: mux}, svc: svc}
}

func (s *Server) Start() error {
	log.Printf("starting sbi on %s", s.cfg.Address)
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.http.Shutdown(ctx)
}

type apiHandler struct { svc *service.Service }

func (h *apiHandler) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req domain.CreateSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, err.Error(), 400); return }
		sess, err := h.svc.Create(req)
		if err != nil { http.Error(w, err.Error(), 400); return }
		writeJSON(w, 201, sess)
	case http.MethodGet:
		writeJSON(w, 200, h.svc.List())
	default:
		w.WriteHeader(405)
	}
}

func (h *apiHandler) handleSessionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/mbs/sessions/"):]
	if id == "" { http.Error(w, "missing id", 400); return }
	switch r.Method {
	case http.MethodGet:
		s, err := h.svc.Get(id); if err != nil { http.Error(w, err.Error(), 404); return }
		writeJSON(w, 200, s)
	case http.MethodPatch:
		var req domain.UpdateSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, err.Error(), 400); return }
		s, err := h.svc.Update(id, req); if err != nil { http.Error(w, err.Error(), 400); return }
		writeJSON(w, 200, s)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		var (
			s   *domain.MBSSession
			err error
		)
		switch action {
		case "activate":
			s, err = h.svc.Activate(id)
		case "deactivate":
			s, err = h.svc.Deactivate(id)
		default:
			http.Error(w, fmt.Sprintf("unknown action %q", action), 400); return
		}
		if err != nil { http.Error(w, err.Error(), 400); return }
		writeJSON(w, 200, s)
	case http.MethodDelete:
		if err := h.svc.Delete(id); err != nil { http.Error(w, err.Error(), 404); return }
		w.WriteHeader(204)
	default:
		w.WriteHeader(405)
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
