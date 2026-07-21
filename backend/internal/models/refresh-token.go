package models

import "time"

type RefreshToken struct {
	UUID      string    `json:"uuid"`
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
