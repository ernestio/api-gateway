package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func authenticate(c echo.Context) error {
	var u User

	username := c.FormValue("username")
	password := c.FormValue("password")

	// Find user, sending the auth request as payload
	req := fmt.Sprintf(`{"username": "%s"}`, username)
	msg, err := n.Request("user.find", []byte(req), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	err = json.Unmarshal(msg.Data, &u)
	if err != nil {
		return ErrBadReqBody
	}

	if u.ID == "" {
		return ErrUnauthorized
	}

	if u.Username == username && u.Password == password {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		token.Claims["username"] = u.Username
		token.Claims["admin"] = u.Admin
		token.Claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return ErrUnauthorized
}
