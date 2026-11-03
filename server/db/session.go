package db

import (
	"sync"
	"time"
)

type SessionRepo interface {
	Insert(string, int, time.Time) (bool, error)
	Select(string) (Session, bool, error)
	Delete(string) (bool, error)
}

type sessionRepo struct {
	sync.RWMutex
	store map[string]Session
}

func newSessionRepo() (SessionRepo, error) {
	return &sessionRepo{
		store: make(map[string]Session),
	}, nil
}

func (h *sessionRepo) Select(token string) (Session, bool, error) {
	h.RLock()
	defer h.RUnlock()
	data, exists := h.store[token]
	if !exists {
		return Session{}, false, nil
	}
	return data, true, nil
}

func (h *sessionRepo) Insert(token string, userId int, expiry time.Time) (bool, error) {
	h.Lock()
	defer h.Unlock()
	h.store[token] = Session{
		Token: token,
		UserId: userId,
		Expiry: expiry,
	}
	return true, nil
}

func (h *sessionRepo) Delete(token string) (bool, error) {
	h.Lock()
	defer h.Unlock()
	delete(h.store, token)
	return true, nil
}
