package web

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

// SessionData holds info about a session
type SessionData struct {
	id     string
	cookie string
}

// SessionManager interface
type SessionManager interface {
	New() (SessionData, error)
	Delete(id string) error
	Get(id string) (SessionData, error)
}

// InMemorySessionManager stores sessions in memory
type InMemorySessionManager struct {
	sessions map[string]SessionData
	mu       sync.Mutex
}

// NewInMemorySessionManager creates a new manager
func NewInMemorySessionManager() *InMemorySessionManager {
	return &InMemorySessionManager{
		sessions: make(map[string]SessionData),
	}
}

// New creates a new session
func (m *InMemorySessionManager) New() (SessionData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()
	cookie := uuid.NewString()
	s := SessionData{
		ID:     id,
		Cookie: cookie,
	}
	m.sessions[id] = s
	return s, nil
}

// Delete removes a session
func (m *InMemorySessionManager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[id]; !ok {
		return errors.New("session not found")
	}
	delete(m.sessions, id)
	return nil
}

// Get retrieves a session by ID
func (m *InMemorySessionManager) Get(id string) (SessionData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[id]
	if !ok {
		return SessionData{}, errors.New("session not found")
	}
	return s, nil
}
