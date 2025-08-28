package domain

import (
	"time"
)

type MBSSession struct {
	ID            string            
	ServiceArea   string            
	Tsi           string                    // Temporary Session Identifier
	ContentDesc   string            
	Policies      map[string]string 
	CreatedAt     time.Time         
	UpdatedAt     time.Time         
	State         SessionState      
}

type SessionState string

const (
	SessionStateCreated   SessionState = "CREATED"
	SessionStateActive    SessionState = "ACTIVE"
	SessionStateInactive  SessionState = "INACTIVE"
	SessionStateDeleted   SessionState = "DELETED"
)

// Requests inspired by 29.532 but simplified

type CreateSessionRequest struct {
	ServiceArea string            
	ContentDesc string            
	Policies    map[string]string 
}

type UpdateSessionRequest struct {
	Policies map[string]string 
	State    *SessionState     
}
