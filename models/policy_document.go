/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// PolicyDocument holds the policy revision response from policy store
type PolicyDocument struct {
	ID         int    `json:"id"`
	PolicyID   int    `json:"policy_id"`
	Revision   int    `json:"revision"`
	Username   string `json:"username"`
	Definition string `json:"definition"`
	CreatedAt  string `json:"created_at"`
}

// Validate : validates the policy
func (p *PolicyDocument) Validate() error {
	if p.Definition == "" {
		return errors.New("Policy definition is empty")
	}

	return nil
}

// Map : maps a policy document revision from a request's body and validates the input
func (p *PolicyDocument) Map(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindAll : Searches for all policys on the system
func (p *PolicyDocument) FindAll(documents *[]PolicyDocument) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel("policy").FindBy(query, documents)
}

// FindByID : Gets a policy revision by policy name and id
func (p *PolicyDocument) FindByID(name, id string, documents *[]PolicyDocument) (err error) {
	query := make(map[string]interface{})

	query["policy"] = name
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	return NewBaseModel("policy_document").FindBy(query, documents)
}

// FindByID : Gets a policy revision by policy name and id
func (p *PolicyDocument) FindByPolicy(name string, documents *[]PolicyDocument) (err error) {
	query := make(map[string]interface{})

	query["policy"] = name
	return NewBaseModel("policy_document").FindBy(query, documents)
}

// FindByID : Gets a policy revision by policy name and revision number
func (p *PolicyDocument) GetByRevision(name, revision string, policy *PolicyDocument) (err error) {
	query := make(map[string]interface{})

	query["policy_name"] = name
	if query["revision"], err = strconv.Atoi(revision); err != nil {
		return err
	}
	return NewBaseModel("policy_document").GetBy(query, policy)
}

func (p *PolicyDocument) Save() (err error) {
	return NewBaseModel("policy_document").Save(p)
}
