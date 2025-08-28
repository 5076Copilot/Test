package service

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"mbsmf/internal/config"
	"mbsmf/internal/metrics"
	"mbsmf/internal/models"
	memstore "mbsmf/internal/store/memory"
)

type mockNotifier struct{ events []string }

func (m *mockNotifier) DispatchSubscriptionEvent(sub models.Subscription, event string, payload any) {
	m.events = append(m.events, event)
}

func newTestService() (*MBService, *mockNotifier) {
	store := memstore.NewMemoryStore()
	metrics := metrics.NewRegistry()
	logger := log.New(os.Stdout, "test ", 0)
	nt := &mockNotifier{}
	cfg := config.LoadFromEnv()
	return NewMBService(store, nt, logger, metrics, cfg), nt
}

func TestCreateSession_Valid(t *testing.T) {
	svc, _ := newTestService()
	sess, err := svc.CreateSession(models.CreateSessionRequest{TMGI: "1", ServiceArea: []string{"A"}, QoS: "gold"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sess.ID == "" || sess.TMGI != "1" {
		t.Fatalf("unexpected session: %+v", sess)
	}
	if sess.State != models.SessionStateActive {
		t.Fatalf("expected active state")
	}
}

func TestCreateSession_Invalid(t *testing.T) {
	svc, _ := newTestService()
	_, err := svc.CreateSession(models.CreateSessionRequest{TMGI: "", ServiceArea: nil})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestSubscriptionAndNotify(t *testing.T) {
	svc, nt := newTestService()
	// create subscription for session.create
	sub, err := svc.CreateSubscription(models.CreateSubscriptionRequest{CallbackURI: "http://localhost/cb", EventTypes: []string{"session.create"}})
	if err != nil {
		t.Fatalf("unexpected error creating subscription: %v", err)
	}
	if sub.ID == "" {
		t.Fatalf("expected subscription id")
	}
	// create session -> should trigger event
	_, err = svc.CreateSession(models.CreateSessionRequest{TMGI: "100", ServiceArea: []string{"X"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(nt.events) == 0 || nt.events[0] != "session.create" {
		t.Fatalf("expected session.create event, got %v", nt.events)
	}
}

func TestPolicyCRUD(t *testing.T) {
	svc, _ := newTestService()
	p, err := svc.CreatePolicy(models.CreatePolicyRequest{Rules: []string{"allow:tmgi=1"}})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if p.ID == "" {
		t.Fatalf("expected id")
	}
	b, _ := json.Marshal(svc.ListPolicies())
	if len(b) == 0 {
		t.Fatalf("expected list json")
	}
	if err := svc.DeletePolicy(p.ID); err != nil {
		t.Fatalf("delete policy error: %v", err)
	}
}

