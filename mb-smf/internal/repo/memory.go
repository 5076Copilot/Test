package repo

import (
	"errors"
	"sync"
	"time"

	"example.com/mb-smf/internal/domain"
)

type MemorySessionRepo struct {
	mu       sync.RWMutex
	sessions map[string]*domain.MBSSession
}

func NewMemorySessionRepo() *MemorySessionRepo {
	return &MemorySessionRepo{sessions: make(map[string]*domain.MBSSession)}
}

func (r *MemorySessionRepo) Create(sess *domain.MBSSession) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.sessions[sess.ID]; ok { return errors.New("exists") }
	sess.CreatedAt = time.Now().UTC()
	sess.UpdatedAt = sess.CreatedAt
	r.sessions[sess.ID] = sess
	return nil
}

func (r *MemorySessionRepo) Get(id string) (*domain.MBSSession, error) {
	r.mu.RLock(); defer r.mu.RUnlock()
	s, ok := r.sessions[id]
	if !ok { return nil, errors.New("not found") }
	return s, nil
}

func (r *MemorySessionRepo) Update(id string, mut func(*domain.MBSSession) error) (*domain.MBSSession, error) {
	r.mu.Lock(); defer r.mu.Unlock()
	s, ok := r.sessions[id]
	if !ok { return nil, errors.New("not found") }
	if err := mut(s); err != nil { return nil, err }
	s.UpdatedAt = time.Now().UTC()
	return s, nil
}

func (r *MemorySessionRepo) List() []*domain.MBSSession {
	r.mu.RLock(); defer r.mu.RUnlock()
	out := make([]*domain.MBSSession, 0, len(r.sessions))
	for _, s := range r.sessions { out = append(out, s) }
	return out
}

func (r *MemorySessionRepo) Delete(id string) error {
	r.mu.Lock(); defer r.mu.Unlock()
	if _, ok := r.sessions[id]; !ok { return errors.New("not found") }
	delete(r.sessions, id)
	return nil
}
