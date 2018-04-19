/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package controllers

import (
	"github.com/ernestio/api-gateway/controllers/policies"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/labstack/echo"
)

// GetPoliciesHandler : responds to GET /policies/ with a list of all
// policies
func GetPoliciesHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)

	st, b := policies.List(au)

	return h.Respond(c, st, b)
}

// GetPolicyHandler : responds to GET /policies/:id:/ with the specified
// user details
func GetPolicyHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	name := c.Param("policy")

	st, b := policies.Get(au, name)

	return h.Respond(c, st, b)
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

// GetPolicyDocumentsHandler : responds to GET /policies/ with a list of all
// policies
func GetPolicyDocumentsHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	name := c.Param("policy")

	st, b := policies.ListDocuments(au, name)

	return h.Respond(c, st, b)
}

// GetPolicyDocumentHandler : responds to GET /policies/:id:/revisions/:rev:
func GetPolicyDocumentHandler(c echo.Context) error {
	au := AuthenticatedUser(c)
	name := c.Param("policy")
	revision := c.Param("revision")

	st, b := policies.GetDocument(au, name, revision)

	return h.Respond(c, st, b)
}

// CreatePolicyDocumentHandler : responds to POST /policies/ by creating a policy
// on the data store
func CreatePolicyDocumentHandler(c echo.Context) (err error) {
	au := AuthenticatedUser(c)
	name := c.Param("policy")

	st, b := h.IsAuthorized(&au, "policys/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = policies.CreateDocument(au, name, body)
	}

	return h.Respond(c, st, b)
}
