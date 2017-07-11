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
	return genericList(c, "notification", notifications.List)
}

// CreateNotificationHandler : responds to POST /notifications/ by creating a notification
// on the data store
func CreateNotificationHandler(c echo.Context) (err error) {
	return genericCreate(c, "notification", notifications.Create)
}

// DeleteNotificationHandler : responds to DELETE /notifications/:id: by deleting an
// existing notification
func DeleteNotificationHandler(c echo.Context) (err error) {
	return genericDelete(c, "notification", notifications.Delete)
}

// UpdateNotificationHandler : ...
func UpdateNotificationHandler(c echo.Context) (err error) {
	return genericUpdate(c, "notification", notifications.Update)
}

// AddServiceToNotificationHandler : ...
func AddServiceToNotificationHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "notifications/add_service")
	if st != 200 {
		return c.JSONBlob(st, b)
	}

	id := c.Param("notification")
	service := c.Param("service")

	st, b = notifications.AddService(au, id, service)

	return c.JSONBlob(st, b)
}

// RmServiceToNotificationHandler : ...
func RmServiceToNotificationHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, "notifications/rm_service")
	if st == 200 {
		return c.JSONBlob(st, b)
	}

	id := c.Param("notification")
	service := c.Param("service")

	st, b = notifications.RmService(au, id, service)

	return c.JSONBlob(st, b)
}
