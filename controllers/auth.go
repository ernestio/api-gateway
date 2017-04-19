/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// Secret : TODO
var Secret string

// AuthenticatedUser : Get the authenticated user from the JWT Token
func AuthenticatedUser(c echo.Context) models.User {
	var u models.User

	user := c.Get("user").(*jwt.Token)

	claims, ok := user.Claims.(jwt.MapClaims)
	if ok {
		u.Username = claims["username"].(string)
		u.GroupID = int(claims["group_id"].(float64))
		u.Admin = claims["admin"].(bool)
	}

	return u
}

// AuthenticateHandler : manages user authentication
func AuthenticateHandler(c echo.Context) error {
	var u models.User

	username := c.FormValue("username")
	password := c.FormValue("password")

	if err := u.FindByUserName(username, &u); err != nil {
		return h.ErrGatewayTimeout
	}

	if u.ID == 0 {
		return h.ErrUnauthorized
	}

	if u.Username == username && u.ValidPassword(password) {
		claims := make(jwt.MapClaims)

		claims["group_id"] = u.GroupID
		claims["username"] = u.Username
		claims["admin"] = u.Admin
		claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(Secret))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return h.ErrUnauthorized
}
