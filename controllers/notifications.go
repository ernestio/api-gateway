/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

// GetNotificationsHandler : responds to GET /notifications/ with a list of all
// notifications
func GetNotificationsHandler(c echo.Context) (err error) {
	var notifications []models.Notification
	var body []byte
	var notification models.Notification

	if err := Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	if err = notification.FindAll(&notifications); err != nil {
		return err
	}

	if body, err = json.Marshal(notifications); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// CreateNotificationHandler : responds to POST /notifications/ by creating a notification
// on the data store
func CreateNotificationHandler(c echo.Context) (err error) {
	var l models.Notification
	var body []byte

	if err := Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if l.Map(data) != nil {
		return h.ErrBadReqBody
	}

	if err = l.Save(); err != nil {
		return c.JSONBlob(400, []byte(err.Error()))
	}

	if body, err = json.Marshal(l); err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, body)
}

// DeleteNotificationHandler : responds to DELETE /notifications/:id: by deleting an
// existing notification
func DeleteNotificationHandler(c echo.Context) (err error) {
	var l models.Notification

	if err := Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if l.Map(data) != nil {
		return h.ErrBadReqBody
	}

	if err := l.Delete(); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}

// Licensed : Checks if the current api is running with premium support
func Licensed() error {
	if len(os.Getenv("ERNEST_PREMIUM")) == 0 {
		return errors.New("You're running ernest community edition, please contact R3Labs for premium support")
	}
	return nil
}

// UpdateNotificationHandler : ...
func UpdateNotificationHandler(c echo.Context) (err error) {
	var d models.Notification
	var existing models.Notification
	var body []byte

	if err := Licensed(); err != nil {
		return c.JSONBlob(403, []byte(err.Error()))
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return h.ErrBadReqBody
	}

	if d.Map(data) != nil {
		return h.ErrBadReqBody
	}

	id := c.Param("notification")
	if err = existing.FindByID(id, &existing); err != nil {
		return err
	}

	existing.Config = d.Config

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return h.ErrInternal
	}

	if body, err = json.Marshal(d); err != nil {
		return h.ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// AddServiceToNotificationHandler : ...
func AddServiceToNotificationHandler(c echo.Context) (err error) {
	var d models.Notification
	var existing models.Notification
	var body []byte

	if err := Licensed(); err != nil {
		return err
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	id := c.Param("notification")
	if err = existing.FindByID(id, &existing); err != nil {
		return err
	}

	if body, err = json.Marshal(d); err != nil {
		return h.ErrInternal
	}

	service := c.Param("service")
	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m == service {
			return c.JSONBlob(http.StatusOK, body)
		}
		if m != "" {
			newMembers = append(newMembers, m)
		}
	}

	members = append(newMembers, service)
	existing.Members = strings.Join(members, ",")

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return h.ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}

// RmServiceToNotificationHandler : ...
func RmServiceToNotificationHandler(c echo.Context) (err error) {
	var d models.Notification
	var existing models.Notification
	var body []byte

	if err := Licensed(); err != nil {
		return err
	}

	au := AuthenticatedUser(c)
	if au.Admin == false {
		return h.ErrUnauthorized
	}

	id := c.Param("notification")
	if err = existing.FindByID(id, &existing); err != nil {
		return err
	}

	service := c.Param("service")
	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m != service {
			newMembers = append(newMembers, m)
		}
	}
	existing.Members = strings.Join(newMembers, ",")

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return h.ErrInternal
	}

	if body, err = json.Marshal(d); err != nil {
		return h.ErrInternal
	}

	return c.JSONBlob(http.StatusOK, body)
}
