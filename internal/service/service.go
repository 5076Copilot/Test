package service

import (
	"errors"
	"log"
	"net/url"
	"time"

	"mbsmf/internal/config"
	"mbsmf/internal/metrics"
	"mbsmf/internal/models"
	memstore "mbsmf/internal/store/memory"
)

var (
	ErrInvalid = errors.New("invalid input")
)

type Notifier interface {
	DispatchSubscriptionEvent(sub models.Subscription, event string, payload any)
}

type MBService struct {
	store    *memstore.MemoryStore
	notify   Notifier
	logger   *log.Logger
	metrics  *metrics.Registry
	config   config.Config
}

func NewMBService(store *memstore.MemoryStore, notify Notifier, logger *log.Logger, metrics *metrics.Registry, cfg config.Config) *MBService {
	return &MBService{store: store, notify: notify, logger: logger, metrics: metrics, config: cfg}
}

// Sessions API
func (s *MBService) CreateSession(req models.CreateSessionRequest) (models.Session, error) {
	if req.TMGI == "" || len(req.ServiceArea) == 0 {
		return models.Session{}, ErrInvalid
	}
	sess := models.Session{
		ID:          models.NewID(),
		TMGI:        req.TMGI,
		ServiceArea: append([]string(nil), req.ServiceArea...),
		QoS:         req.QoS,
		State:       models.SessionStateActive,
	}
	sess = s.store.CreateSession(sess)
	s.metrics.SessionsActive.Add(1)
	// Notify subscriptions interested in session.create
	for _, sub := range s.store.ListSubscriptions() {
		if contains(sub.EventTypes, "session.create") && (len(sub.FilterTMGI) == 0 || contains(sub.FilterTMGI, sess.TMGI)) {
			s.notify.DispatchSubscriptionEvent(sub, "session.create", sess)
		}
	}
	return sess, nil
}

func (s *MBService) DeleteSession(id string) error {
	if err := s.store.DeleteSession(id); err != nil {
		return err
	}
	s.metrics.SessionsActive.Add(-1)
	for _, sub := range s.store.ListSubscriptions() {
		if contains(sub.EventTypes, "session.delete") {
			s.notify.DispatchSubscriptionEvent(sub, "session.delete", map[string]string{"id": id})
		}
	}
	return nil
}

func (s *MBService) ListSessions() []models.Session { return s.store.ListSessions() }

// Subscriptions API
func (s *MBService) CreateSubscription(req models.CreateSubscriptionRequest) (models.Subscription, error) {
	if _, err := url.ParseRequestURI(req.CallbackURI); err != nil {
		return models.Subscription{}, ErrInvalid
	}
	if len(req.EventTypes) == 0 {
		return models.Subscription{}, ErrInvalid
	}
	sub := models.Subscription{
		ID:          models.NewID(),
		CallbackURI: req.CallbackURI,
		EventTypes:  append([]string(nil), req.EventTypes...),
		FilterTMGI:  append([]string(nil), req.FilterTMGI...),
		CreatedAt:   time.Now().UTC(),
	}
	sub = s.store.CreateSubscription(sub)
	s.metrics.SubscriptionsActive.Add(1)
	return sub, nil
}

func (s *MBService) DeleteSubscription(id string) error {
	if err := s.store.DeleteSubscription(id); err != nil {
		return err
	}
	s.metrics.SubscriptionsActive.Add(-1)
	return nil
}

func (s *MBService) ListSubscriptions() []models.Subscription { return s.store.ListSubscriptions() }

// Policies API
func (s *MBService) CreatePolicy(req models.CreatePolicyRequest) (models.Policy, error) {
	if len(req.Rules) == 0 {
		return models.Policy{}, ErrInvalid
	}
	p := models.Policy{ID: models.NewID(), Rules: append([]string(nil), req.Rules...)}
	p = s.store.CreatePolicy(p)
	return p, nil
}

func (s *MBService) DeletePolicy(id string) error { return s.store.DeletePolicy(id) }
func (s *MBService) ListPolicies() []models.Policy { return s.store.ListPolicies() }

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

