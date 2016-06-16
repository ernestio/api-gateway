/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

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
	msg, err := n.Request("user.get", []byte(req), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if responseErr(msg) != nil {
		return ErrUnauthorized
	}

	err = json.Unmarshal(msg.Data, &u)
	if err != nil {
		return ErrInternal
	}

	if u.ID == 0 {
		return ErrUnauthorized
	}

	if u.Username == username && u.ValidPassword(password) {
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
