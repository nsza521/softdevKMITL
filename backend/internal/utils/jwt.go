package utils

import (
	"os"
	"time"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type BlacklistTokens struct {
	sync.RWMutex
	tokens map[string]int64
}
var blacklistedTokens = BlacklistTokens{
	tokens: make(map[string]int64),
}

func GenerateJWTToken(userID uuid.UUID, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Must change later
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func BlacklistToken(token string, expiry int64) {
	blacklistedTokens.Lock()
	defer blacklistedTokens.Unlock()
	blacklistedTokens.tokens[token] = expiry
}

func IsTokenBlacklisted(token string) bool {

	blacklistedTokens.RLock()
	expiry, exists := blacklistedTokens.tokens[token]
	blacklistedTokens.RUnlock()

	if !exists {
		return false
	}

	// Remove token from blacklist if expired
	if time.Now().Unix() > expiry {
		blacklistedTokens.Lock()
		delete(blacklistedTokens.tokens, token)
		blacklistedTokens.Unlock()
		return false
	}

	return true
}

func BlacklistCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			blacklistedTokens.Lock()
			now := time.Now().Unix()
			for token, expiry := range blacklistedTokens.tokens {
				if now > expiry {
					delete(blacklistedTokens.tokens, token)
				}
			}
			blacklistedTokens.Unlock()
		}
	}()
}
