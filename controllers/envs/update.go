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

// Update : responds to PUT /projects/:project:/envs/:env/ by updating an
// existing environment
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var resp []byte
	var e models.Env
	var input models.Env
	var p models.Project
	var r models.Role
	var roles []models.Role
	var pRoles []models.Role

	computedRoles := make(map[string]models.Role, 0)

	if err = json.Unmarshal(body, &input); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get existing environment
	if err = e.FindByName(name); err != nil {
		return 404, models.NewJSONError(err.Error())
	}

	if err = p.FindByID(e.ProjectID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 404, models.NewJSONError("Specified environment name does not exist")
		}
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal error")
	}

	if err = r.FindAllByResource(e.Project, p.GetType(), &pRoles); err == nil {
		for _, v := range pRoles {
			computedRoles[v.UserID] = v
		}
	}
	if err = r.FindAllByResource(e.GetID(), e.GetType(), &roles); err == nil {
		for _, v := range roles {
			computedRoles[v.UserID] = v
		}
	}

	for _, v := range computedRoles {
		e.Members = append(e.Members, v)
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, input.GetType(), name); st != 200 {
		return st, res
	}

	e.Options = input.Options
	e.Schedules = input.Schedules
	e.Credentials = input.Credentials

	if err = e.Save(); err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	for _, ir := range input.Members {
		// create role
		if ir.ID == 0 {
			if !au.IsAdmin() {
				if ok := au.IsOwner(ir.ResourceType, ir.ResourceID); !ok {
					return 403, models.NewJSONError("You're not authorized to perform this action")
				}
			}

			err = ir.Save()
			if err != nil {
				h.L.Error(err.Error())
				return http.StatusBadRequest, models.NewJSONError(err.Error())
			}

			continue
		}

		for _, er := range e.Members {
			// update role
			if ir.ID == er.ID && ir.Role != er.Role {
				if strings.Contains(er.ResourceID, "/") {
					return http.StatusBadRequest, models.NewJSONError("project memberships must be modified on the project")
				}

				if !au.IsAdmin() {
					if ok := au.IsOwner(ir.ResourceType, ir.ResourceID); !ok {
						return 403, models.NewJSONError("You're not authorized to perform this action")
					}
				}

				err = ir.Save()
				if err != nil {
					h.L.Error(err.Error())
					return http.StatusBadRequest, models.NewJSONError(err.Error())
				}
			}
		}
	}

	for _, er := range e.Members {
		var exists bool

		for _, ir := range input.Members {
			if ir.ID == er.ID {
				exists = true
			}
		}

		// delete roles
		if !exists {
			if strings.Contains(er.ResourceID, "/") {
				return http.StatusBadRequest, models.NewJSONError("project memberships must be removed on the project")
			}

			if !au.IsAdmin() {
				if ok := au.IsOwner(er.ResourceType, er.ResourceID); !ok {
					return 403, models.NewJSONError("You're not authorized to perform this action")
				}
			}

			err = er.Delete()
			if err != nil {
				h.L.Error(err.Error())
				return http.StatusBadRequest, models.NewJSONError(err.Error())
			}
		}
	}

	e.Members = input.Members

	resp, err = json.Marshal(e)
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}

	return http.StatusOK, resp
}
