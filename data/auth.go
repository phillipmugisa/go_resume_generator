package data

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

const Session_duration = time.Hour * 24 * 3 // 3 days

type Session struct {
	User       User
	Key        string
	Expires_on time.Time
	Expired    bool
}

func (u User) NewSession() (*Session, error) {
	key, err := generateSessionKey(30)
	if err != nil {
		return nil, err
	}
	return &Session{
		User:       u,
		Key:        key,
		Expires_on: time.Now().Add(Session_duration),
		Expired:    false,
	}, nil
}

func generateSessionKey(length int) (string, error) {
	// Generate a random sequence of bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes into a base64 string
	sessionKey := base64.URLEncoding.EncodeToString(randomBytes)

	return sessionKey, nil
}
