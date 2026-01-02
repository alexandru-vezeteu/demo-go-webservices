package domain

import "time"

type TokenBlacklist interface {
	Add(token string, reason string, expiresAt time.Time) error
	IsBlacklisted(token string) (bool, string)
	Remove(token string) error
}
