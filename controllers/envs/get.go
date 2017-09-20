/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Get : responds to GET /services/:service with the
// details of an existing service
func Get(au models.User, name string) (int, []byte) {
	var o views.BuildRender
	var err error
	var body []byte
	var e models.Env
	var r models.Role
	var p models.Project
	var roles []models.Role

	if err = e.FindByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if e.ID == 0 {
		return 404, []byte("Specified environment name does not exist")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), name); st != 200 {
		return st, res
	}

	if err := r.FindAllByResource(e.GetID(), e.GetType(), &roles); err == nil {
		for _, v := range roles {
			e.Roles = append(e.Roles, v.UserID+" ("+v.Role+")")
		}
	}

	if err := p.FindByID(int(e.ProjectID)); err != nil {
		return 404, []byte("Project not found")
	}

	e.Project = p.Name
	e.Provider = p.Type

	if body, err = o.ToJSON(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
