/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/labstack/echo"
)

// Logger holds the logger response from logger
type Logger struct {
	Type        string `json:"type"`
	Logfile     string `json:"logfile"`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
	Timeout     int    `json:"timeout"`
	Token       string `json:"token"`
	Environment string `json:"environment"`
}

// Validate : validates the logger
func (l *Logger) Validate() error {
	if l.Type == "" {
		return errors.New("Logger type is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (l *Logger) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &l)
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

// FindAll : Searches for all loggers on the system
func (l *Logger) FindAll(loggers *[]Logger) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("logger").FindBy(query, loggers); err != nil {
		return err
	}
	return nil
}

// Save : calls logger.set with the marshalled current logger
func (l *Logger) Save() (err error) {
	if err := NewBaseModel("logger").Save(l); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a logger by its type
func (l *Logger) Delete() (err error) {
	query := make(map[string]interface{})
	query["type"] = l.Type
	if err := NewBaseModel("logger").Delete(query); err != nil {
		return err
	}
	return nil
}
