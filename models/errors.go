/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// Error : the default error type for responses
type Error struct {
	Message    string                 `json:"message"`
	Validation *BuildValidateResponse `json:"validation,omitempty"`
}

// ToJSON : marshals error to json
func (e Error) ToJSON() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		h.L.WithFields(logrus.Fields{
			"input": e.Message,
		}).Error("Couldn't marshal given input")
	}
	return data
}

// NewJSONError : constructs a json payload
func NewJSONError(message string) []byte {
	return Error{Message: message}.ToJSON()
}

// NewJSONValidationError : constructs a json payload
func NewJSONValidationError(message string, validation *BuildValidateResponse) []byte {
	return Error{message, validation}.ToJSON()
}
