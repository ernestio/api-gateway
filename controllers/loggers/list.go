/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package loggers

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /loggers/ with a list of all
// loggers
func List(au models.User) (int, []byte) {
	var err error
	var loggers []models.Logger
	var body []byte
	var logger models.Logger

	if err = logger.FindAll(&loggers); err != nil {
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(loggers); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
