/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/policies"
	"github.com/labstack/echo"
)

// GetPoliciesHandler : responds to GET /policies/ with a list of all
// policies
func GetPoliciesHandler(c echo.Context) (err error) {
	return genericList(c, "policy", policies.List)
}

// GetPolicyHandler : responds to GET /policies/:id:/ with the specified
// user details
func GetPolicyHandler(c echo.Context) error {
	return genericGet(c, "policy", policies.Get)
}

// CreatePolicyHandler : responds to POST /policies/ by creating a policy
// on the data store
func CreatePolicyHandler(c echo.Context) (err error) {
	return genericCreate(c, "policy", policies.Create)
}

// DeletePolicyHandler : responds to DELETE /policies/:id: by deleting an
// existing policy
func DeletePolicyHandler(c echo.Context) (err error) {
	return genericDelete(c, "policy", policies.Delete)
}

// UpdatePolicyHandler : ...
func UpdatePolicyHandler(c echo.Context) (err error) {
	return genericUpdate(c, "policy", policies.Update)
}
