/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /notifications/ by creating a notification
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var l models.Notification
	var err error

	if l.Map(body) != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err = l.Save(); err != nil {
		return 400, []byte(err.Error())
	}

	if body, err = json.Marshal(l); err != nil {
		return 500, []byte("Internal error")
	}
	return http.StatusOK, body
}
