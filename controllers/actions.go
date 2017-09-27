/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"strings"

	"github.com/ernestio/api-gateway/controllers/builds"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// ActionHandler : handles different actions that can be triggered on an env
func ActionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "services/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	action, err := mapAction(c)
	if err != nil {
		return h.Respond(c, 400, []byte(err.Error()))
	}

	switch action.Type {
	case "import":
		filters := strings.Split(c.QueryParams().Get("filters"), ",")
		st, b = builds.Import(au, buildID(c), filters)
	default:
		return h.Respond(c, 400, []byte("unsupported action"))
	}

	return h.Respond(c, st, b)
}
