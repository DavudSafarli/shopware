package inmem

import (
	"context"
	"redirectware/internal"
	"sync"
)

// Storage is an in-memory implementation for tests.
type Storage struct {
	mu         sync.RWMutex
	rules      map[string]internal.FullMatchRule // key: FromCanonical
	welcomeURL string
}

func New() *Storage {
	return &Storage{
		rules:      make(map[string]internal.FullMatchRule),
		welcomeURL: "",
	}
}

func (s *Storage) AddFullMatchRule(ctx context.Context, rule *internal.FullMatchRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rule == nil {
		return nil
	}
	s.rules[rule.FromCanonical] = *rule
	return nil
}

func (s *Storage) FindFullMatchRule(ctx context.Context, canonicalPath string) (rule internal.FullMatchRule, ok bool, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, exists := s.rules[canonicalPath]
	if !exists {
		return internal.FullMatchRule{}, false, nil
	}
	return r, true, nil
}

func (s *Storage) GetWelcomePageURL(ctx context.Context) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.welcomeURL, nil
}

func (s *Storage) SetWelcomePageURL(ctx context.Context, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.welcomeURL = url
	return nil
}
