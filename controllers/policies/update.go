/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : ...
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var d models.Policy
	var existing models.Policy

	if d.Map(body) != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err = existing.FindByName(name, &existing); err != nil {
		return 404, []byte("Not found")
	}

	existing.Definition = d.Definition
	existing.Environments = d.Environments

	if len(existing.Environments) == 0 {
		existing.Environments = make([]string, 0)
	}

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(existing); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
