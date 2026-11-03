package service

import (
	"time"
	"crypto/rand"
	"encoding/base64"
	"splitify/db"

	"golang.org/x/crypto/bcrypt"
)

const EXPIRY_DURATION = 1 * time.Minute

type AuthService interface {
	Register(username string, password string) (bool, error)
	Login(username string, password string) (bool, string, error)
	Verify(token string) (bool, int, error)
}

type authService struct {
	userRepo db.UserRepo
	sessionRepo db.SessionRepo
}

func NewAuthService(loginRepo db.UserRepo, sessionRepo db.SessionRepo) AuthService {
	return &authService{
		userRepo: loginRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *authService) Register(username string, password string) (bool, error) {
	passwordBytes := []byte(password)

	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	_, err = h.userRepo.Insert(username, string(hashed))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *authService) Login(username string, password string) (bool, string, error) {
	actual, err := h.userRepo.SelectByUsername(username)
	if err != nil {
		return false, "", err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(actual.Password), []byte(password))
	if err != nil {
		return false, "", nil
	}

	token, err := getToken()

	success, err := h.sessionRepo.Insert(token, actual.Id, time.Now().Add(EXPIRY_DURATION))
	if err != nil {
		return false, "", nil
	}

	return success, token, nil
}

func getToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	return token, nil
}

func (h *authService) Verify(token string) (bool, int, error) {
	sess, found, err := h.sessionRepo.Select(token)
	if err != nil {
		return false, -1, err
	}

	if !found {
		return false, -1, nil
	}

	if time.Now().After(sess.Expiry) {
		h.sessionRepo.Delete(token)
		return false, -1, nil
	}

	return true, sess.UserId, nil
}
