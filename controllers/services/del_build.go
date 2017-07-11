/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package services

import (
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// DelBuild : will delete the specified build from a service
func DelBuild(au models.User, id string) (int, []byte) {
	if id == "" {
		h.L.Debug("Empty id")
		return 400, []byte("Invalid build id")
	}

	build, err := au.GetBuild(id)
	if err != nil {
		return 404, []byte("Not found")
	}

	if err := build.Delete(); err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Oops something went wrong")
	}

	return 200, []byte("Build successfully removed")
}
