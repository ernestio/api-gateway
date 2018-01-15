/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /policies/:name: by deleting an
// existing policy
func Delete(au models.User, name string) (int, []byte) {
	var err error
	var existing models.Policy

	if err = existing.GetByName(name, &existing); err != nil {
		return 404, []byte("policy not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.DeletePolicy, existing.GetType(), existing.GetID()); st != 200 {
		return st, res
	}

	if err := existing.Delete(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("policy deleted")
}
