/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /policies/:name: by deleting an
// existing policy
func Delete(au models.User, name string) (int, []byte) {
	var err error
	var existing models.Policy

	if err = existing.FindByName(name, &existing); err != nil {
		return 404, []byte("Not found")
	}

	if err := existing.Delete(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("Policy deleted")
}
