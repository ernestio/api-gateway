/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"io/ioutil"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetUsersHandler : responds to GET /users/ with a list of all
// users for admin, and all users in your group for other
// users
func GetUsersHandler(c echo.Context) error {
	var users []models.User

	au := AuthenticatedUser(c)
	if err := au.FindAll(&users); err != nil {
		return err
	}

	for i := 0; i < len(users); i++ {
		users[i].Redact()
		users[i].Improve()
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserHandler : responds to GET /users/:id:/ with the specified
// user details
func GetUserHandler(c echo.Context) error {
	var user models.User

	au := AuthenticatedUser(c)
	if err := au.FindByID(c.Param("user"), &user); err != nil {
		return err
	}
	user.Redact()

	return c.JSON(http.StatusOK, user)
}

// CreateUserHandler : responds to POST /users/ by creating a user
// on the data store
func CreateUserHandler(c echo.Context) error {
	var u models.User
	var existing models.User

	if AuthenticatedUser(c).Admin != true {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}
	if u.Map(data) != nil {
		return h.ErrBadReqBody
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

// UpdateUserHandler : responds to PUT /users/:id: by updating an existing
// user
func UpdateUserHandler(c echo.Context) error {
	var u models.User
	var existing models.User

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}
	if u.Map(data) != nil {
		return h.ErrBadReqBody
	}

	// Check if authenticated user is admin or updating itself
	au := AuthenticatedUser(c)
	if au.Username != u.Username && au.Admin != true {
		return h.ErrUnauthorized
	}

	// Check user exists
	if err := au.FindByID(c.Param("user"), &existing); err != nil {
		return err
	}

	if existing.ID == 0 {
		return h.ErrNotFound
	}

	// Check a non-admin user is not trying to change their group
	if au.Admin != true && u.GroupID != existing.GroupID {
		return h.ErrUnauthorized
	}

	// Check the old password if it is present
	if u.OldPassword != "" && !existing.ValidPassword(u.OldPassword) {
		return h.ErrUnauthorized
	}

	if err := u.Save(); err != nil {
		return err
	}

	u.Redact()

	return c.JSON(http.StatusOK, u)
}

// DeleteUserHandler : responds to DELETE /users/:id: by deleting an
// existing user
func DeleteUserHandler(c echo.Context) error {
	var au models.User

	if au = AuthenticatedUser(c); au.Admin != true {
		return h.ErrUnauthorized
	}

	if err := au.Delete(c.Param("user")); err != nil {
		return err
	}

	return c.String(http.StatusOK, "User successfully deleted")
}
