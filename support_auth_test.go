package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateTestToken(group int, user string, admin bool) *jwt.Token {
	// Set claims
	claims := ErnestClaims{
		GroupID:  group,
		Username: user,
		Admin:    admin,
	}

	claims.ExpiresAt = time.Now().Add(time.Hour * 48).Unix()

	// Create token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
