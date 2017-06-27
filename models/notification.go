/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Sirupsen/logrus"
	h "github.com/ernestio/api-gateway/helpers"
)

// Notification holds the notification response from notification
type Notification struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Config  string `json:"config"`
	Members string `json:"members"`
}

// Validate : validates the notification
func (l *Notification) Validate() error {
	if l.Type == "" {
		return errors.New("Notification type is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (l *Notification) Map(data []byte) error {
	if err := json.Unmarshal(data, &l); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindAll : Searches for all notifications on the system
func (l *Notification) FindAll(notifications *[]Notification) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("notification").FindBy(query, notifications); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a notification by ID
func (l *Notification) FindByID(id string, notification *Notification) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	if err := NewBaseModel("notification").GetBy(query, notification); err != nil {
		return err
	}
	return nil

}

// FindByName : Searches for all notifications with a name equal to the specified
func (l *Notification) FindByName(name string, notification *Notification) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("notification").GetBy(query, notification); err != nil {
		return err
	}
	return nil
}

// Save : calls notification.set with the marshalled current notification
func (l *Notification) Save() (err error) {
	if err := NewBaseModel("notification").Save(l); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a notification by its type
func (l *Notification) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = l.ID
	if err := NewBaseModel("notification").Delete(query); err != nil {
		return err
	}
	return nil
}
