package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Registry struct {
	RequestsTotal         atomic.Int64
	RequestsErrorsTotal   atomic.Int64
	NotificationsSent     atomic.Int64
	NotificationsFailed   atomic.Int64
	SessionsActive        atomic.Int64
	SubscriptionsActive   atomic.Int64
	MbsServicesActive     atomic.Int64
	MbsSessionsActive     atomic.Int64
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintf(w, "# HELP mbsmf_requests_total Total HTTP requests\n")
		fmt.Fprintf(w, "# TYPE mbsmf_requests_total counter\nmbsmf_requests_total %d\n", r.RequestsTotal.Load())
		fmt.Fprintf(w, "# HELP mbsmf_request_errors_total Total HTTP request errors\n")
		fmt.Fprintf(w, "# TYPE mbsmf_request_errors_total counter\nmbsmf_request_errors_total %d\n", r.RequestsErrorsTotal.Load())
		fmt.Fprintf(w, "# HELP mbsmf_notifications_sent_total Total notifications sent\n")
		fmt.Fprintf(w, "# TYPE mbsmf_notifications_sent_total counter\nmbsmf_notifications_sent_total %d\n", r.NotificationsSent.Load())
		fmt.Fprintf(w, "# HELP mbsmf_notifications_failed_total Total notifications failed\n")
		fmt.Fprintf(w, "# TYPE mbsmf_notifications_failed_total counter\nmbsmf_notifications_failed_total %d\n", r.NotificationsFailed.Load())
		fmt.Fprintf(w, "# HELP mbsmf_sessions_active Active sessions\n")
		fmt.Fprintf(w, "# TYPE mbsmf_sessions_active gauge\nmbsmf_sessions_active %d\n", r.SessionsActive.Load())
		fmt.Fprintf(w, "# HELP mbsmf_subscriptions_active Active subscriptions\n")
		fmt.Fprintf(w, "# TYPE mbsmf_subscriptions_active gauge\nmbsmf_subscriptions_active %d\n", r.SubscriptionsActive.Load())
		fmt.Fprintf(w, "# HELP mbsmf_mbs_services_active Active MBS services\n")
		fmt.Fprintf(w, "# TYPE mbsmf_mbs_services_active gauge\nmbsmf_mbs_services_active %d\n", r.MbsServicesActive.Load())
		fmt.Fprintf(w, "# HELP mbsmf_mbs_sessions_active Active MBS sessions\n")
		fmt.Fprintf(w, "# TYPE mbsmf_mbs_sessions_active gauge\nmbsmf_mbs_sessions_active %d\n", r.MbsSessionsActive.Load())
	})
}

