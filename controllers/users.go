/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/users"
	"github.com/labstack/echo"
)

// GetUsersHandler : responds to GET /users/ with a list of all
// users for admin, and all users in your group for other
// users
func GetUsersHandler(c echo.Context) error {
	return genericList(c, "user", users.List)
}

// GetUserHandler : responds to GET /users/:id:/ with the specified
// user details
func GetUserHandler(c echo.Context) error {
	return genericGet(c, "user", users.Get)
}

// CreateUserHandler : responds to POST /users/ by creating a user
// on the data store
func CreateUserHandler(c echo.Context) error {
	return genericCreate(c, "user", users.Create)
}

// UpdateUserHandler : responds to PUT /users/:id: by updating an existing
// user
func UpdateUserHandler(c echo.Context) error {
	return genericUpdate(c, "user", users.Update)
}

// DeleteUserHandler : responds to DELETE /users/:id: by deleting an
// existing user
func DeleteUserHandler(c echo.Context) error {
	return genericDelete(c, "user", users.Delete)
}
