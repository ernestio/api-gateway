/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /notifications/:id: by deleting an
// existing notification
func Delete(au models.User, id string) (int, []byte) {
	var err error
	var existing models.Notification

	if err = existing.FindByID(id, &existing); err != nil {
		return 500, []byte("Internal server error")
	}

	if err := existing.Delete(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("Notification deleted")
}
