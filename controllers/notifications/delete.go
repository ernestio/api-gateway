/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /notifications/:name: by deleting an
// existing notification
func Delete(au models.User, name string) (int, []byte) {
	var err error
	var existing models.Notification

	if !models.IsAlphaNumeric(name) {
		return 404, models.NewJSONError("Notification name contains invalid characters")
	}

	if err = existing.FindByName(name, &existing); err != nil {
		return 404, models.NewJSONError("Not found")
	}

	if err := existing.Delete(); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, []byte(`{"status": "Notification deleted"}`)
}
