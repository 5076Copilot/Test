package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type SessionState string

const (
	SessionStateActive   SessionState = "ACTIVE"
	SessionStateInactive SessionState = "INACTIVE"
)

type Session struct {
	ID           string       `json:"id"`
	TMGI         string       `json:"tmgi"`
	ServiceArea  []string     `json:"serviceArea"`
	QoS          string       `json:"qos"`
	State        SessionState `json:"state"`
	CreatedAt    time.Time    `json:"createdAt"`
	ModifiedAt   time.Time    `json:"modifiedAt"`
}

type CreateSessionRequest struct {
	TMGI        string   `json:"tmgi"`
	ServiceArea []string `json:"serviceArea"`
	QoS         string   `json:"qos"`
}

type Subscription struct {
	ID          string    `json:"id"`
	CallbackURI string    `json:"callbackUri"`
	EventTypes  []string  `json:"eventTypes"`
	FilterTMGI  []string  `json:"filterTmgi"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateSubscriptionRequest struct {
	CallbackURI string   `json:"callbackUri"`
	EventTypes  []string `json:"eventTypes"`
	FilterTMGI  []string `json:"filterTmgi"`
}

type Policy struct {
	ID        string    `json:"id"`
	Rules     []string  `json:"rules"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreatePolicyRequest struct {
	Rules []string `json:"rules"`
}

func NewID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

