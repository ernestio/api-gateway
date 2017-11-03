/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/builds"
	"github.com/ernestio/api-gateway/controllers/envs"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// ActionHandler : handles different actions that can be triggered on an env
func ActionHandler(c echo.Context) error {
	au := AuthenticatedUser(c)

	action, err := mapAction(c)
	if err != nil {
		return h.Respond(c, 400, []byte(err.Error()))
	}

	st, b := h.IsAuthorized(&au, "envs/"+action.Type)
	if st != 200 {
		return h.Respond(c, st, b)
	}

	switch action.Type {
	case "import":
		st, b = builds.Import(au, envName(c), action)
	case "reset":
		st, b = envs.Reset(au, envName(c), action)
	case "sync":
		st, b = envs.Sync(au, envName(c), action)
	case "resolve":
		st, b = envs.Resolve(au, envName(c), action)
	case "review":
		st, b = builds.Approval(au, envName(c), action)
	default:
		return h.Respond(c, 400, []byte("unsupported action"))
	}

	return h.Respond(c, st, b)
}
