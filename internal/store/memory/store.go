package memory

import (
	"errors"
	"sync"
	"time"

	"mbsmf/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type MemoryStore struct {
	mu            sync.RWMutex
	sessions      map[string]models.Session
	subscriptions map[string]models.Subscription
	policies      map[string]models.Policy
	mbsServices   map[string]models.MbsService
	mbsSessions   map[string]models.MbsSession
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions:      make(map[string]models.Session),
		subscriptions: make(map[string]models.Subscription),
		policies:      make(map[string]models.Policy),
		mbsServices:   make(map[string]models.MbsService),
		mbsSessions:   make(map[string]models.MbsSession),
	}
}

// Sessions
func (s *MemoryStore) CreateSession(sess models.Session) models.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess.CreatedAt = time.Now().UTC()
	sess.ModifiedAt = sess.CreatedAt
	s.sessions[sess.ID] = sess
	return sess
}

func (s *MemoryStore) GetSession(id string) (models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.sessions[id]; ok {
		return v, nil
	}
	return models.Session{}, ErrNotFound
}

func (s *MemoryStore) DeleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[id]; !ok {
		return ErrNotFound
	}
	delete(s.sessions, id)
	return nil
}

func (s *MemoryStore) ListSessions() []models.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.Session, 0, len(s.sessions))
	for _, v := range s.sessions {
		res = append(res, v)
	}
	return res
}

// Subscriptions
func (s *MemoryStore) CreateSubscription(sub models.Subscription) models.Subscription {
	s.mu.Lock()
	defer s.mu.Unlock()
	sub.CreatedAt = time.Now().UTC()
	s.subscriptions[sub.ID] = sub
	return sub
}

func (s *MemoryStore) GetSubscription(id string) (models.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.subscriptions[id]; ok {
		return v, nil
	}
	return models.Subscription{}, ErrNotFound
}

func (s *MemoryStore) DeleteSubscription(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[id]; !ok {
		return ErrNotFound
	}
	delete(s.subscriptions, id)
	return nil
}

func (s *MemoryStore) ListSubscriptions() []models.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.Subscription, 0, len(s.subscriptions))
	for _, v := range s.subscriptions {
		res = append(res, v)
	}
	return res
}

// Policies
func (s *MemoryStore) CreatePolicy(p models.Policy) models.Policy {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.CreatedAt = time.Now().UTC()
	s.policies[p.ID] = p
	return p
}

func (s *MemoryStore) GetPolicy(id string) (models.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.policies[id]; ok {
		return v, nil
	}
	return models.Policy{}, ErrNotFound
}

func (s *MemoryStore) DeletePolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.policies[id]; !ok {
		return ErrNotFound
	}
	delete(s.policies, id)
	return nil
}

func (s *MemoryStore) ListPolicies() []models.Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.Policy, 0, len(s.policies))
	for _, v := range s.policies {
		res = append(res, v)
	}
	return res
}

// MBS Services
func (s *MemoryStore) CreateMbsService(m models.MbsService) models.MbsService {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.CreatedAt = time.Now().UTC()
	s.mbsServices[m.ID] = m
	return m
}

func (s *MemoryStore) GetMbsService(id string) (models.MbsService, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.mbsServices[id]; ok {
		return v, nil
	}
	return models.MbsService{}, ErrNotFound
}

func (s *MemoryStore) DeleteMbsService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mbsServices[id]; !ok {
		return ErrNotFound
	}
	delete(s.mbsServices, id)
	return nil
}

func (s *MemoryStore) ListMbsServices() []models.MbsService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.MbsService, 0, len(s.mbsServices))
	for _, v := range s.mbsServices {
		res = append(res, v)
	}
	return res
}

// MBS Sessions
func (s *MemoryStore) CreateMbsSession(m models.MbsSession) models.MbsSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.CreatedAt = time.Now().UTC()
	m.ModifiedAt = m.CreatedAt
	s.mbsSessions[m.ID] = m
	return m
}

func (s *MemoryStore) GetMbsSession(id string) (models.MbsSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.mbsSessions[id]; ok {
		return v, nil
	}
	return models.MbsSession{}, ErrNotFound
}

func (s *MemoryStore) DeleteMbsSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mbsSessions[id]; !ok {
		return ErrNotFound
	}
	delete(s.mbsSessions, id)
	return nil
}

func (s *MemoryStore) ListMbsSessions() []models.MbsSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.MbsSession, 0, len(s.mbsSessions))
	for _, v := range s.mbsSessions {
		res = append(res, v)
	}
	return res
}

