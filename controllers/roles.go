/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/roles"
	h "github.com/ernestio/api-gateway/helpers"

	"github.com/labstack/echo"
)

// GetRolesHandler : responds to GET /roles/ with a list of all
// roles for admin, and all roles in your group for other
// roles
func GetRolesHandler(c echo.Context) error {
	return genericList(c, "role", roles.List)
}

// GetRoleHandler : responds to GET /roles/:role/ with the role
// details
func GetRoleHandler(c echo.Context) error {
	return genericGet(c, "role", roles.Get)
}

// CreateRoleHandler : responds to POST /roles/ by creating a
// role on the data store
func CreateRoleHandler(c echo.Context) (err error) {
	return genericCreate(c, "role", roles.Create)
}

// DeleteRoleByIDHandler : responds to DELETE /roles/:id: by deleting an
// existing role
func DeleteRoleByIDHandler(c echo.Context) error {
	return genericDelete(c, "role", roles.DeleteByID)
}

// DeleteRoleHandler : responds to DELETE /roles/:id: by deleting an
// existing role
func DeleteRoleHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "roles/delete")
	if st == 200 {
		st = 500
		b = []byte("Invalid input")
		body, err := h.GetRequestBody(c)
		if err == nil {
			st, b = roles.Delete(au, body)
		}
	}

	return h.Respond(c, st, b)
}
