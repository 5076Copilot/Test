package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"mbsmf/internal/config"
	"mbsmf/internal/metrics"
	"mbsmf/internal/models"
	"mbsmf/internal/service"
)

type router struct {
	svc     *service.MBService
	logger  *log.Logger
	metrics *metrics.Registry
	cfg     config.Config
}

func NewRouter(svc *service.MBService, logger *log.Logger, metrics *metrics.Registry, cfg config.Config) http.Handler {
	r := &router{svc: svc, logger: logger, metrics: metrics, cfg: cfg}
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", r.handleHealth)
	mux.Handle("/metrics", r.metrics.Handler())

	mux.HandleFunc("/nmbsf/v1/sessions", r.sessions)
	mux.HandleFunc("/nmbsf/v1/sessions/", r.sessionByID)

	mux.HandleFunc("/nmbsf/v1/subscriptions", r.subscriptions)
	mux.HandleFunc("/nmbsf/v1/subscriptions/", r.subscriptionByID)

	mux.HandleFunc("/nmbsf/v1/policies", r.policies)
	mux.HandleFunc("/nmbsf/v1/policies/", r.policyByID)

	// Rel-18 MBS resources
	mux.HandleFunc("/nmbsf/v1/mbs-services", r.mbsServices)
	mux.HandleFunc("/nmbsf/v1/mbs-services/", r.mbsServiceByID)
	mux.HandleFunc("/nmbsf/v1/mbs-sessions", r.mbsSessions)
	mux.HandleFunc("/nmbsf/v1/mbs-sessions/", r.mbsSessionByID)

	return withCommon(r.logger, r.metrics, mux)
}

func withCommon(logger *log.Logger, metrics *metrics.Registry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestsTotal.Add(1)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (r *router) handleHealth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.cfg.ReadinessDelay > 0 {
		// simple readiness gate after startup
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Sessions
func (r *router) sessions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var body models.CreateSessionRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		res, err := r.svc.CreateSession(body)
		if err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, res)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, r.svc.ListSessions())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *router) sessionByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := req.URL.Path[len("/nmbsf/v1/sessions/"):]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.svc.DeleteSession(id); err != nil {
		writeErr(r.logger, r.metrics, w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Subscriptions
func (r *router) subscriptions(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var body models.CreateSubscriptionRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		res, err := r.svc.CreateSubscription(body)
		if err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, res)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, r.svc.ListSubscriptions())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *router) subscriptionByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := req.URL.Path[len("/nmbsf/v1/subscriptions/"):]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.svc.DeleteSubscription(id); err != nil {
		writeErr(r.logger, r.metrics, w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Policies
func (r *router) policies(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var body models.CreatePolicyRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		res, err := r.svc.CreatePolicy(body)
		if err != nil {
			writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, res)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, r.svc.ListPolicies())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *router) policyByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := req.URL.Path[len("/nmbsf/v1/policies/"):]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := r.svc.DeletePolicy(id); err != nil {
		writeErr(r.logger, r.metrics, w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MBS Services
func (r *router) mbsServices(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        var body models.CreateMbsServiceRequest
        if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
            writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
            return
        }
        res, err := r.svc.CreateMbsService(body)
        if err != nil {
            writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
            return
        }
        writeJSON(w, http.StatusCreated, res)
    case http.MethodGet:
        writeJSON(w, http.StatusOK, r.svc.ListMbsServices())
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (r *router) mbsServiceByID(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    id := req.URL.Path[len("/nmbsf/v1/mbs-services/"):]
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if err := r.svc.DeleteMbsService(id); err != nil {
        writeErr(r.logger, r.metrics, w, http.StatusNotFound, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// MBS Sessions
func (r *router) mbsSessions(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        var body models.CreateMbsSessionRequest
        if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
            writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
            return
        }
        res, err := r.svc.CreateMbsSession(body)
        if err != nil {
            writeErr(r.logger, r.metrics, w, http.StatusBadRequest, err)
            return
        }
        writeJSON(w, http.StatusCreated, res)
    case http.MethodGet:
        writeJSON(w, http.StatusOK, r.svc.ListMbsSessions())
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (r *router) mbsSessionByID(w http.ResponseWriter, req *http.Request) {
    if req.Method != http.MethodDelete {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    id := req.URL.Path[len("/nmbsf/v1/mbs-sessions/"):]
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if err := r.svc.DeleteMbsSession(id); err != nil {
        writeErr(r.logger, r.metrics, w, http.StatusNotFound, err)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func writeErr(logger *log.Logger, m *metrics.Registry, w http.ResponseWriter, code int, err error) {
	m.RequestsErrorsTotal.Add(1)
	logger.Printf("http error %d: %v", code, err)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

