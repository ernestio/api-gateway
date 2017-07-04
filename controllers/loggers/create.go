/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package loggers

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /loggers/ by creating a logger
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var l models.Logger
	var err error

	if au.Admin == false {
		return 403, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action")
	}

	if l.Map(body) != nil {
		return 400, []byte("Invalid input")
	}

	if err = l.Save(); err != nil {
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(l); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
