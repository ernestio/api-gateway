/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ernestio/api-gateway/models"
)

// AddService : ...
func AddService(au models.User, id, service string) (int, []byte) {
	var err error
	var d models.Notification
	var existing models.Notification
	var body []byte

	if au.Admin == false {
		return 403, []byte("You should provide admin credentials to perform this action")
	}

	if err = existing.FindByID(id, &existing); err != nil {
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m == service {
			return http.StatusOK, body
		}
		if m != "" {
			newMembers = append(newMembers, m)
		}
	}

	members = append(newMembers, service)
	existing.Members = strings.Join(members, ",")

	if err = existing.Save(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
