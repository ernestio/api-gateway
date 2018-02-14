/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"net/http"
	"strings"

	"github.com/ernestio/api-gateway/models"
)

// AddEnv : ...
func AddEnv(au models.User, name, env string) (int, []byte) {
	var err error
	var existing models.Notification
	var body []byte

	if err = existing.FindByName(name, &existing); err != nil {
		return 500, models.NewJSONError("Noification not found")
	}

	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m == env {
			return http.StatusOK, body
		}
		if m != "" {
			newMembers = append(newMembers, m)
		}
	}

	members = append(newMembers, env)
	existing.Members = strings.Join(members, ",")

	if err = existing.Save(); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
