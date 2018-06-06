/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"encoding/json"
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /services/:service with the
// details of an existing service
func Get(au models.User, name string) (int, []byte) {
	var err error
	var body []byte
	var e models.Env
	var p models.Project
	var r models.Role
	var roles []models.Role
	var pRoles []models.Role

	if err = e.FindByName(name); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 404, models.NewJSONError("Specified environment name does not exist")
		}
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal error")
	}

	if err = p.FindByID(e.ProjectID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 404, models.NewJSONError("Specified environment name does not exist")
		}
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal error")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), name); st != 200 {
		return st, res
	}

	computedRoles := make(map[string]models.Role, 0)
	if err := r.FindAllByResource(e.GetProject(), p.GetType(), &pRoles); err == nil {
		for _, v := range pRoles {
			computedRoles[v.UserID] = v
		}
	}
	if err := r.FindAllByResource(e.GetID(), e.GetType(), &roles); err == nil {
		for _, v := range roles {
			computedRoles[v.UserID] = v
		}
	}

	for _, v := range computedRoles {
		e.Members = append(e.Members, v)
	}

	err = e.Redact()
	if err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	if body, err = json.Marshal(e); err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	return http.StatusOK, body
}
