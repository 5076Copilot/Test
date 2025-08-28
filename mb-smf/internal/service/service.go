package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"example.com/mb-smf/internal/domain"
	"example.com/mb-smf/internal/pfcp"
)

type SessionRepo interface {
	Create(*domain.MBSSession) error
	Get(string) (*domain.MBSSession, error)
	Update(string, func(*domain.MBSSession) error) (*domain.MBSSession, error)
	List() []*domain.MBSSession
	Delete(string) error
}

type Service struct {
	repo SessionRepo
	pfcp *pfcp.Client
}

func New(repo SessionRepo, pfcpClient *pfcp.Client) *Service {
	return &Service{repo: repo, pfcp: pfcpClient}
}

func (s *Service) Create(req domain.CreateSessionRequest) (*domain.MBSSession, error) {
	id := randomID()
	sess := &domain.MBSSession{ID: id, ServiceArea: req.ServiceArea, ContentDesc: req.ContentDesc, Policies: req.Policies, State: domain.SessionStateCreated}
	if err := s.repo.Create(sess); err != nil { return nil, err }
	// send PFCP Session Establishment (placeholder body)
	_ = s.pfcp.SetupSession([]byte("establish"))
	return sess, nil
}

func (s *Service) Activate(id string) (*domain.MBSSession, error) {
	return s.repo.Update(id, func(m *domain.MBSSession) error {
		if m.State == domain.SessionStateDeleted { return errors.New("deleted") }
		m.State = domain.SessionStateActive
		_ = s.pfcp.ModifySession([]byte("activate"))
		return nil
	})
}

func (s *Service) Deactivate(id string) (*domain.MBSSession, error) {
	return s.repo.Update(id, func(m *domain.MBSSession) error {
		if m.State == domain.SessionStateDeleted { return errors.New("deleted") }
		m.State = domain.SessionStateInactive
		_ = s.pfcp.ModifySession([]byte("deactivate"))
		return nil
	})
}

func (s *Service) Update(id string, req domain.UpdateSessionRequest) (*domain.MBSSession, error) {
	return s.repo.Update(id, func(m *domain.MBSSession) error {
		if req.Policies != nil { m.Policies = req.Policies }
		if req.State != nil { m.State = *req.State }
		_ = s.pfcp.ModifySession([]byte("update"))
		return nil
	})
}

func (s *Service) Get(id string) (*domain.MBSSession, error) { return s.repo.Get(id) }
func (s *Service) List() []*domain.MBSSession { return s.repo.List() }

func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil { return err }
	_ = s.pfcp.DeleteSession([]byte("delete"))
	return nil
}

func randomID() string {
	var b [8]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
