/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /notifications/ with a list of all
// notifications
func List(au models.User) (int, []byte) {
	var err error
	var notifications []models.Notification
	var body []byte
	var notification models.Notification

	if err = notification.FindAll(&notifications); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal serveier error")
	}

	if body, err = json.Marshal(notifications); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal serveier error")
	}
	return http.StatusOK, body
}
