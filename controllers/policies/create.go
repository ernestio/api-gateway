/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /policies/ by creating a policy
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var l models.Policy
	var existing models.Policy
	var err error

	if l.Map(body) != nil {
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	if err = l.GetByName(l.Name, &existing); err == nil {
		return 409, models.NewJSONError("policy already exists")
	}

	if err = l.Save(); err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	if err := au.SetOwner(&l); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	if body, err = json.Marshal(l); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}
	return http.StatusOK, body
}
