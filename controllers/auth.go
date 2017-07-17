/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/dgrijalva/jwt-go"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

const RequiredCliVersion string = "2.2.0"

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

	// validate input
	var vu models.User
	vu.Username = c.FormValue("username")
	vu.Password = c.FormValue("password")

	if err := vu.Validate(); err != nil {
		h.L.Error(err.Error())
		return echo.NewHTTPError(400, err.Error())
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if err := u.FindByUserName(username, &u); err != nil {
		return echo.NewHTTPError(403, "The keypair user / password does not match any user on the database, please try again")
	}

	if u.ID == 0 {
		return h.ErrUnauthorized
	}

	if u.Username == username && u.ValidPassword(password) {
		// Check CLI version requirement
		req := c.Request()
		uagent := req.Header["User-Agent"]

		for _, v := range uagent {
			if strings.Contains(v, "Ernest/") {
				parts := strings.Split(v, "/")
				agentVersion := parts[1]

				rv, err := semver.Make(RequiredCliVersion)
				if err != nil {
					return err
				}
				ev, err := semver.Make(agentVersion)
				if err != nil {
					return err
				}
				if ev.LT(rv) {
					err := fmt.Sprintf("Ernest CLI %s is not supported by this server.\nPlease upgrade http://docs.ernest.io/downloads/", agentVersion)
					return echo.NewHTTPError(403, err)
				}
			}
		}

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

	return echo.NewHTTPError(403, "The keypair user / password does not match any user on the database, please try again")
}
