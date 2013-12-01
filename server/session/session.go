package session

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"
)

// SessionManager manages sessions that are valid for a certain time.
type SessionManager struct {
	l         sync.RWMutex
	sessions  map[string]*session
	timeLimit time.Duration
}

// NewSessionManager will expire sessions after them not being used for
// timeLimit.
func NewSessionManager(timeLimit time.Duration) *SessionManager {
	return &SessionManager{
		l:         sync.RWMutex{},
		sessions:  make(map[string]*session),
		timeLimit: timeLimit,
	}
}

// NewSession ties a username to a remoteAddr, and returns a token to use with
// to authenticate the session in the future.
func (s *SessionManager) NewSession(username, remoteAddr string) (string, error) {
	session, err := newSession(username, remoteAddr)
	if err != nil {
		return "", err
	}
	s.l.Lock()
	s.sessions[session.token] = session
	s.l.Unlock()
	return session.token, nil
}

// UseSession takes a token and a remoteAddr string, verify they match known
// records and return the matching username.  If the token/remoteAddr pair don't
// match with records, the return values are "" and false.
func (s *SessionManager) UseSession(token, remoteAddr string) (string, bool) {
	s.l.RLock()
	session, ok := s.sessions[token]
	s.l.RUnlock()
	if !ok {
		return "", false
	}

	session.l.RLock()
	if time.Since(session.lastAccess) > s.timeLimit {
		session.l.RUnlock()
		// revoke token
		s.l.Lock()
		delete(s.sessions, token)
		s.l.Unlock()
		// and don't authorize
		return "", false
	}

	if session.token != token {
		session.l.RUnlock()
		return "", false
	}

	if session.remoteAddr != remoteAddr {
		session.l.RUnlock()
		return "", false
	}
	username := session.username
	session.l.RUnlock()
	session.l.Lock()
	session.lastAccess = time.Now()
	session.l.Unlock()
	return username, true
}

// DeleteSession removes the session associated with the token, if any
func (s *SessionManager) DeleteSession(token string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.sessions, token)
}

// CollectGarbage will remove all the tokens that are expired.
func (s *SessionManager) CollectGarbage() {
	go func() {
		now := time.Now()
		s.l.RLock()
		for token, session := range s.sessions {
			session.l.RLock()
			if now.Sub(session.lastAccess) > s.timeLimit {
				s.l.RUnlock()
				s.l.Lock()
				delete(s.sessions, token)
				s.l.Unlock()
				s.l.RLock()
			}
			session.l.RUnlock()
		}
		s.l.RUnlock()
	}()
}

type session struct {
	l sync.RWMutex
	// For exhaustion
	lastAccess time.Time
	// All two should match
	token      string
	remoteAddr string // request should come from unique IP (can be spoofed)
	// to retrieve the username from a token
	username string
}

func newSession(username, remoteAddr string) (*session, error) {
	token, err := computeToken()
	if err != nil {
		return nil, fmt.Errorf("couldn't compute token, %v", err)
	}
	return &session{
		l:          sync.RWMutex{},
		lastAccess: time.Now(),
		token:      token,
		remoteAddr: remoteAddr,
		username:   username,
	}, nil
}

func computeToken() (string, error) {
	hash := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, hash)
	if err != nil {
		return "", fmt.Errorf("reading random, %v", err)
	} else if n != len(hash) {
		return "", fmt.Errorf("wanted %d random bytes but got %d", len(hash), n)
	}
	return url.QueryEscape(string(hash)), nil
}
