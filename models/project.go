/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	aes "github.com/ernestio/crypto/aes"
	"github.com/sirupsen/logrus"
)

// Project holds the project response from datacenter-store
type Project struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	Credentials  map[string]interface{} `json:"credentials,omitempty"`
	Environments []string               `json:"environments,omitempty"`
	Members      []Role                 `json:"members,omitempty"`
}

// Validate the project
func (d *Project) Validate() error {
	if d.Name == "" {
		return errors.New("Project name is empty")
	}

	if !IsAlphaNumeric(d.Name) {
		return errors.New("Project Name contains invalid characters")
	}

	if strings.Contains(d.Name, EnvNameSeparator) {
		return errors.New("Project name does not support char '" + EnvNameSeparator + "' as part of its name")
	}

	if d.Type == "" {
		return errors.New("Project type is empty")
	}

	switch d.Type {
	case "aws", "azure", "vcloud":
		return nil
	default:
		return errors.New("Project type is not one of the following: 'aws', 'azure' or 'vcloud'")
	}
}

// Map : maps a project from a request's body and validates the input
func (d *Project) Map(data []byte) error {
	if err := json.Unmarshal(data, &d); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindByName : Searches for all projects with a name equal to the specified
func (d *Project) FindByName(name string) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel(d.getStore()).GetBy(query, d)
}

// FindByID : Gets a model by its id
func (d *Project) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	return NewBaseModel(d.getStore()).GetBy(query, d)
}

// FindByIDs : Gets a model by its id
func (d *Project) FindByIDs(ids []string, ds *[]Project) (err error) {
	query := make(map[string]interface{})
	query["names"] = ids
	return NewBaseModel(d.getStore()).FindBy(query, ds)
}

// FindAll : Searches for all entities on the store current user
// has access to
func (d *Project) FindAll(au User, projects *[]Project) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel(d.getStore()).FindBy(query, projects)
}

// Save : calls datacenter.set with the marshalled current entity
func (d *Project) Save() (err error) {
	return NewBaseModel(d.getStore()).Save(d)
}

// Delete : will delete a project by its id
func (d *Project) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = d.ID
	return NewBaseModel(d.getStore()).Delete(query)
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (d *Project) Redact() {
	for k := range d.Credentials {
		if k != "region" && k != "vdc" && k != "username" && k != "vcloud_url" {
			delete(d.Credentials, k)
		}
	}
}

// Improve : adds extra data to this entity
func (d *Project) Improve() {
}

// Envs : Get the envs related with current project
func (d *Project) Envs() (envs []Env, err error) {
	var s Env
	err = s.FindByProjectID(d.ID, &envs)
	return
}

// GetID : ID getter
func (d *Project) GetID() string {
	return d.Name
}

// GetType : Gets the resource type
func (d *Project) GetType() string {
	return "project"
}

// Override : override not empty parameters with the given project ones
func (d *Project) Override(dt Project) {
	for k, v := range dt.Credentials {
		d.Credentials[k] = v
	}
}

// Encrypt : encrypts sensible data
func (d *Project) Encrypt() {
	for k, v := range d.Credentials {
		xc, ok := v.(string)
		if !ok {
			continue
		}

		d.Credentials[k], _ = crypt(xc)
	}
}

func decrypt(s string) (string, error) {
	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	if s != "" {
		encrypted, err := crypto.Decrypt(s, key)
		if err != nil {
			return "", err
		}
		s = encrypted
	}

	return s, nil
}

func crypt(s string) (string, error) {
	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	if s != "" {
		encrypted, err := crypto.Encrypt(s, key)
		if err != nil {
			return "", err
		}
		s = encrypted
	}

	return s, nil
}

func (d *Project) getStore() string {
	return "datacenter"
}

// IsAzure : check if the current project is of type azure
func (d *Project) IsAzure() bool {
	return (d.Type == "azure" || d.Type == "azure-fake")
}
