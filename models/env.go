/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// EnvNameSeparator : environment name separator
var EnvNameSeparator = "/"

// Env holds the environment response from service-store
type Env struct {
	ID          int                    `json:"id"`
	ProjectID   int                    `json:"project_id"`
	Project     string                 `json:"project,omitempty"`
	Provider    string                 `json:"provider,omitempty"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Roles       []string               `json:"roles,omitempty"`
}

// Validate the env
func (s *Env) Validate() error {
	if s.Name == "" {
		return errors.New("Environment name is empty")
	}

	return nil
}

// Map : maps a env from a request's body and validates the input
func (s *Env) Map(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := s.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Warning("Input is not valid")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// Find : Searches for all envs with filters
func (s *Env) Find(query map[string]interface{}, envs *[]Env) (err error) {
	return NewBaseModel(s.getStore()).FindBy(query, envs)
}

// FindByName : Searches for all envs with a name equal to the specified
func (s *Env) FindByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel(s.getStore()).GetBy(query, s)
}

// FindByID : Gets a model by its id
func (s *Env) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	return NewBaseModel(s.getStore()).GetBy(query, s)
}

// FindAll : Searches for all envs s on the store current user
// has access to
func (s *Env) FindAll(au User, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel(s.getStore()).FindBy(query, envs)
}

// Save : calls env.set with the marshalled
func (s *Env) Save() (err error) {
	return NewBaseModel(s.getStore()).Save(s)
}

// Delete : will delete a env by its id
func (s *Env) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = s.ID
	return NewBaseModel(s.getStore()).Delete(query)
}

// DeleteByName : will delete a env by its name
func (s *Env) DeleteByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel(s.getStore()).Delete(query)
}

// FindByProjectID : find a envs for the given project id
func (s *Env) FindByProjectID(id int, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	query["project_id"] = id
	return NewBaseModel(s.getStore()).FindBy(query, envs)
}

/*

// RequestSync : calls service.sync with the given raw message
func (s *Env) RequestSync() error {
	if err := N.Publish(s.getStore()+".sync", []byte(`{"id":"`+s.ID+`"}`)); err != nil {
		h.L.Error(err.Error())
		return err
	}
	return nil
}

*/

// GetID : ID getter
func (s *Env) GetID() string {
	return s.Name
}

// GetType : Gets the resource type
func (s *Env) GetType() string {
	return "environment"
}

// getStore : Gets the store name
func (s *Env) getStore() string {
	return "environment"
}
