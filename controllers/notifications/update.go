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

// Update : ...
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var d models.Notification
	var existing models.Notification

	if d.Map(body) != nil {
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	err = d.Validate()
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	if err = existing.FindByName(name, &existing); err != nil {
		return 404, models.NewJSONError("Not found")
	}

	existing.Config = d.Config

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
