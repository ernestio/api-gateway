/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// RmEnv : ...
func RmEnv(au models.User, name, env string) (int, []byte) {
	var err error
	var existing models.Notification
	var body []byte

	if err = existing.FindByName(name, &existing); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m != env {
			newMembers = append(newMembers, m)
		}
	}
	existing.Members = strings.Join(newMembers, ",")

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
