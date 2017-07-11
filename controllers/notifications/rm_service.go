/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// RmService : ...
func RmService(au models.User, id, service string) (int, []byte) {
	var err error
	var d models.Notification
	var existing models.Notification
	var body []byte

	if err = existing.FindByID(id, &existing); err != nil {
		return 500, []byte("Internal server error")
	}

	members := strings.Split(existing.Members, ",")
	newMembers := make([]string, 0)
	for _, m := range members {
		if m != service {
			newMembers = append(newMembers, m)
		}
	}
	existing.Members = strings.Join(newMembers, ",")

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
