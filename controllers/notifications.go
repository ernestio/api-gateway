/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/notifications"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetNotificationsHandler : responds to GET /notifications/ with a list of all
// notifications
func GetNotificationsHandler(c echo.Context) (err error) {
	if err := h.Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}
	au := AuthenticatedUser(c)
	s, b := notifications.List(au)

	return c.JSONBlob(s, b)
}

// CreateNotificationHandler : responds to POST /notifications/ by creating a notification
// on the data store
func CreateNotificationHandler(c echo.Context) (err error) {
	if err := h.Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = notifications.Create(au, body)
	}

	return c.JSONBlob(s, b)
}

// DeleteNotificationHandler : responds to DELETE /notifications/:id: by deleting an
// existing notification
func DeleteNotificationHandler(c echo.Context) (err error) {
	if err := h.Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	id := c.Param("notification")

	s, b := notifications.Delete(au, id)

	return c.JSONBlob(s, b)
}

// UpdateNotificationHandler : ...
func UpdateNotificationHandler(c echo.Context) (err error) {
	s := 500
	b := []byte("Invalid input")
	au := AuthenticatedUser(c)
	id := c.Param("notification")

	body, err := h.GetRequestBody(c)
	if err == nil {
		s, b = notifications.Update(au, id, body)
	}

	return c.JSONBlob(s, b)
}

// AddServiceToNotificationHandler : ...
func AddServiceToNotificationHandler(c echo.Context) (err error) {
	if err := h.Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	id := c.Param("notification")
	service := c.Param("service")

	s, b := notifications.AddService(au, id, service)

	return c.JSONBlob(s, b)
}

// RmServiceToNotificationHandler : ...
func RmServiceToNotificationHandler(c echo.Context) (err error) {
	if err := h.Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	id := c.Param("notification")
	service := c.Param("service")

	s, b := notifications.RmService(au, id, service)

	return c.JSONBlob(s, b)
}
