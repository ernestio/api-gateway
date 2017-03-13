package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateTestToken(group int, user string, admin bool) *jwt.Token {
	claims := make(jwt.MapClaims)

	claims["group_id"] = float64(group)
	claims["username"] = user
	claims["admin"] = admin
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	// Create token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
