/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package middleware

import (
	"github.com/ernestio/api-gateway/controllers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
	"strings"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		au := controllers.AuthenticatedUser(c)
		req := c.Request()
		s := strings.Split(req.RemoteAddr, ":")
		ip := s[0]

		if err := models.Log("api-gateway", "method="+req.Method+" path="+req.RequestURI+" remote-ip="+ip, models.LogInfoLevel, au.Username); err != nil {
			return err
		}

		return next(c)
	}
}
