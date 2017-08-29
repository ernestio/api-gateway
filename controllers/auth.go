/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"net/http"

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
		u.Admin = claims["admin"].(bool)
	}

	return u
}

// AuthenticateHandler manages user authentication
func AuthenticateHandler(c echo.Context) error {
	var u models.User

	u.Username = c.FormValue("username")
	u.Password = c.FormValue("password")

	err := u.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return echo.NewHTTPError(400, err.Error())
	}

	res, err := u.Authenticate()
	if err != nil {
		h.L.Error(err.Error())
		return echo.NewHTTPError(400, err.Error())
	}

	if !res.OK {
		h.L.Error(res.Message + " (" + u.Username + ")")
		return echo.NewHTTPError(403, res.Message)
	}

	if err := h.ValidCliVersion(c.Request()); err != nil {
		h.L.Error(err.Error())
		return echo.NewHTTPError(403, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"token": res.Token})
}
