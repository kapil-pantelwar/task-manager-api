package session

import (
	"time"
)

type Session struct {
	SessionID		string `json:"session_id"`
	Uid       int `json:"uid"`
	TokenHash string `json:"tokenhash"`
	ExpiresAt time.Time `json:"expiresat"`
	IssuedAt  time.Time `json:"issuedat"`
}
