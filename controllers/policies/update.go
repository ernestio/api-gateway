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
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	err = d.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}

	if err = existing.GetByName(name, &existing); err != nil {
		return 404, models.NewJSONError("Not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdatePolicy, existing.GetType(), existing.GetID()); st != 200 {
		return st, res
	}

	existing.Environments = d.Environments

	if len(existing.Environments) == 0 {
		existing.Environments = make([]string, 0)
	}

	existing.Username = au.Username

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	if body, err = json.Marshal(existing); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
