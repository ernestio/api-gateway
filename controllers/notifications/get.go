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

// Get : responds to GET /notifications/:name/ with the notification
// details
func Get(au models.User, name string) (int, []byte) {
	var err error
	var body []byte
	var notification models.Notification

	if err = notification.FindByName(name, &notification); err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Notification not found")
	}

	if body, err = json.Marshal(notification); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal serveier error")
	}
	return http.StatusOK, body
}
