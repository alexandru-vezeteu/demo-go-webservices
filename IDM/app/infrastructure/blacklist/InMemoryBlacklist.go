package blacklist

import (
	"fmt"
	"sync"
	"time"
)

type BlacklistedToken struct {
	Token   string
	Reason  string
	AddedAt time.Time
}

type InMemoryBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]*BlacklistedToken
}

func NewInMemoryBlacklist() *InMemoryBlacklist {
	return &InMemoryBlacklist{
		tokens: make(map[string]*BlacklistedToken),
	}
}

func (b *InMemoryBlacklist) Add(token string, reason string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens[token] = &BlacklistedToken{
		Token:   token,
		Reason:  reason,
		AddedAt: time.Now(),
	}

	return nil
}

func (b *InMemoryBlacklist) IsBlacklisted(token string) (bool, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if entry, exists := b.tokens[token]; exists {
		return true, entry.Reason
	}

	return false, ""
}

func (b *InMemoryBlacklist) Remove(token string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.tokens[token]; !exists {
		return fmt.Errorf("token not found in blacklist")
	}

	delete(b.tokens, token)
	return nil
}

func (b *InMemoryBlacklist) Clear() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.tokens = make(map[string]*BlacklistedToken)
	return nil
}
