/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

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
	Schedules   map[string]interface{} `json:"schedules,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Builds      []Build                `json:"builds,omitempty"`
	Members     []Role                 `json:"members,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	UpdatedAt   string                 `json:"updated_at,omitempty"`
}

// Validate the env
func (e *Env) Validate() error {
	if e.Name == "" {
		return errors.New("Environment name is empty")
	}

	if !IsAlphaNumeric(e.Name) {
		return errors.New("Environment name contains invalid characters")
	}

	return nil
}

// Map : maps a env from a request's body and validates the input
func (e *Env) Map(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := e.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Warning("Input is not valid")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// Find : Searches for all envs with filters
func (e *Env) Find(query map[string]interface{}, envs *[]Env) (err error) {
	return NewBaseModel(e.getStore()).FindBy(query, envs)
}

// FindByName : Searches for all envs with a name equal to the specified
func (e *Env) FindByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel(e.getStore()).GetBy(query, e)
}

// FindByID : Gets a model by its id
func (e *Env) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	return NewBaseModel(e.getStore()).GetBy(query, e)
}

// FindAll : Searches for all envs s on the store current user
// has access to
func (e *Env) FindAll(au User, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel(e.getStore()).FindBy(query, envs)
}

// Save : calls env.set with the marshalled
func (e *Env) Save() (err error) {
	return NewBaseModel(e.getStore()).Save(e)
}

// Delete : will delete a env by its id
func (e *Env) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = e.ID
	return NewBaseModel(e.getStore()).Delete(query)
}

// DeleteByName : will delete a env by its name
func (e *Env) DeleteByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel(e.getStore()).Delete(query)
}

// FindByProjectID : find a envs for the given project id
func (e *Env) FindByProjectID(id int, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	query["project_id"] = id
	return NewBaseModel(e.getStore()).FindBy(query, envs)
}

// FindByProjectName : find a envs for the given project name
func (e *Env) FindByProjectName(name string, envs *[]Env) (err error) {
	query := make(map[string]interface{})
	query["project_name"] = name
	return NewBaseModel(e.getStore()).FindBy(query, envs)
}

// RequestSync : calls environment.sync with the given raw message
func (e *Env) RequestSync(au User) (string, error) {
	req := map[string]interface{}{
		"name":     e.Name,
		"user_id":  au.ID,
		"username": au.Username,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := N.Request(e.getStore()+".sync", data, time.Second*5)
	if err != nil {
		h.L.Error(err.Error())
		return "", err
	}

	var r struct {
		ID string `json:"id"`
	}

	err = json.Unmarshal(resp.Data, &r)
	if err != nil {
		return "", err
	}

	return r.ID, nil
}

// RequestResolve : calls environment.resolve with the given raw message
func (e *Env) RequestResolve(au User, resolution string) (string, error) {
	return e.resolution(au, e.getStore()+".resolve", resolution)
}

// RequestReview : calls build.review with the given raw message
func (e *Env) RequestReview(au User, resolution string) (string, error) {
	return e.resolution(au, "build.review", resolution)
}

func (e *Env) resolution(au User, subject, resolution string) (string, error) {
	req := map[string]interface{}{
		"name":       e.Name,
		"user_id":    au.ID,
		"username":   au.Username,
		"resolution": resolution,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := N.Request(subject, data, time.Second*5)
	if err != nil {
		h.L.Error(err.Error())
		return "", err
	}

	var r struct {
		ID     string `json:"id"`
		Status string `json:"ok"`
		Error  string `json:"_error"`
	}

	err = json.Unmarshal(resp.Data, &r)
	if err != nil {
		return "", err
	}

	if r.Error != "" {
		return r.ID, errors.New(r.Error)
	}

	return r.ID, nil
}

// GetID : ID getter
func (e *Env) GetID() string {
	return e.Name
}

// GetProject : returns the environment's project
func (e *Env) GetProject() string {
	return strings.Split(e.Name, "/")[0]
}

// GetType : Gets the resource type
func (e *Env) GetType() string {
	return "environment"
}

// getStore : Gets the store name
func (e *Env) getStore() string {
	return "environment"
}
