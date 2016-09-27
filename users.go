/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// getUsersHandler : responds to GET /users/ with a list of all
// users for admin, and all users in your group for other
// users
func getUsersHandler(c echo.Context) error {
	var users []User

	au := authenticatedUser(c)
	if err := au.FindAll(&users); err != nil {
		return err
	}

	for i := 0; i < len(users); i++ {
		users[i].Redact()
	}

	return c.JSON(http.StatusOK, users)
}

// getUserHandler : responds to GET /users/:id:/ with the specified
// user details
func getUserHandler(c echo.Context) error {
	var user User

	au := authenticatedUser(c)
	if err := au.FindByID(c.Param("user"), &user); err != nil {
		return err
	}
	user.Redact()

	return c.JSON(http.StatusOK, user)
}

// createUserHandler : responds to POST /users/ by creating a user
// on the data store
func createUserHandler(c echo.Context) error {
	var u User
	var existing User

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	if err := existing.FindByUserName(u.Username, &existing); err == nil {
		return echo.NewHTTPError(409, "Specified user already exists")
	}

	if err := u.Save(); err != nil {
		return err
	}

	u.Redact()

	return c.JSON(http.StatusOK, u)
}

// updateUserHandler : responds to PUT /users/:id: by updating an existing
// user
func updateUserHandler(c echo.Context) error {
	var u User
	var existing User

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	// Check if authenticated user is admin or updating itself
	au := authenticatedUser(c)
	if au.Username != u.Username && au.Admin != true {
		return ErrUnauthorized
	}

	// Check user exists
	if err := au.FindByID(c.Param("user"), &existing); err != nil {
		return err
	}

	if existing.ID == 0 {
		return ErrNotFound
	}

	// Check a non-admin user is not trying to change their group
	if au.Admin != true && u.GroupID != existing.GroupID {
		return ErrUnauthorized
	}

	// Check the old password if it is present
	if u.OldPassword != "" && !existing.ValidPassword(u.OldPassword) {
		return ErrUnauthorized
	}

	if err := u.Save(); err != nil {
		return err
	}

	u.Redact()

	return c.JSON(http.StatusOK, u)
}

// deleteUserHandler : responds to DELETE /users/:id: by deleting an
// existing user
func deleteUserHandler(c echo.Context) error {
	var au User

	if au = authenticatedUser(c); au.Admin != true {
		return ErrUnauthorized
	}

	if err := au.Delete(c.Param("user")); err != nil {
		return err
	}

	return c.String(http.StatusOK, "User successfully deleted")
}
