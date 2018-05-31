/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// Notification holds the notification response from notification
type Notification struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Config  string   `json:"config"`
	Sources []string `json:"sources"`
}

// Validate : validates the notification
func (n *Notification) Validate() error {
	if n.Name == "" {
		return errors.New("Notification must specify a name")
	}

	if !IsAlphaNumeric(n.Name) {
		return errors.New("Notification name contains invalid characters")
	}

	if n.Type == "" {
		return errors.New("Notification type is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (n *Notification) Map(data []byte) error {
	if err := json.Unmarshal(data, &n); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindAll : Searches for all notifications on the system
func (n *Notification) FindAll(notifications *[]Notification) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel("notification").FindBy(query, notifications)
}

// FindByID : Gets a notification by ID
func (n *Notification) FindByID(id string, notification *Notification) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	return NewBaseModel("notification").GetBy(query, notification)
}

// FindByName : Searches for all notifications with a name equal to the specified
func (n *Notification) FindByName(name string, notification *Notification) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel("notification").GetBy(query, notification)
}

// Save : calls notification.set with the marshalled current notification
func (n *Notification) Save() (err error) {
	return NewBaseModel("notification").Save(n)
}

// Delete : will delete a notification by its type
func (n *Notification) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = n.ID
	return NewBaseModel("notification").Delete(query)
}
