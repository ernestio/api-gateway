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

// Datacenter holds the datacenter response from datacenter-store
type Datacenter struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Region          string   `json:"region"`
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	VCloudURL       string   `json:"vcloud_url"`
	VseURL          string   `json:"vse_url"`
	ExternalNetwork string   `json:"external_network"`
	AccessKeyID     string   `json:"aws_access_key_id,omitempty"`
	SecretAccessKey string   `json:"aws_secret_access_key,omitempty"`
	SubscriptionID  string   `json:"azure_subscription_id,omitempty"`
	ClientID        string   `json:"azure_client_id,omitempty"`
	ClientSecret    string   `json:"azure_client_secret,omitempty"`
	TenantID        string   `json:"azure_tenant_id"`
	Environment     string   `json:"azure_environment"`
	Environments    []string `json:"environments,omitempty"`
	Roles           []string `json:"roles,omitempty"`
}

// Validate the datacenter
func (d *Datacenter) Validate() error {
	if d.Name == "" {
		return errors.New("Project name is empty")
	}

	if strings.Contains(d.Name, EnvNameSeparator) {
		return errors.New("Project name does not support char '" + EnvNameSeparator + "' as part of its name")
	}

	if d.Type == "" {
		return errors.New("Project type is empty")
	}

	if d.Username == "" && d.Type != "azure" && d.Type != "azure-fake" {
		return errors.New("Project username is empty")
	}

	if d.Type == "vcloud" && d.VCloudURL == "" {
		return errors.New("Project vcloud url is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (d *Datacenter) Map(data []byte) error {
	if err := json.Unmarshal(data, &d); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindByName : Searches for all datacenters with a name equal to the specified
func (d *Datacenter) FindByName(name string, datacenter *Datacenter) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("datacenter").GetBy(query, datacenter); err != nil {
		return err
	}
	return nil
}

// FindByID : Gets a model by its id
func (d *Datacenter) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel("datacenter").GetBy(query, d); err != nil {
		return err
	}
	return nil
}

// FindByIDs : Gets a model by its id
func (d *Datacenter) FindByIDs(ids []string, ds *[]Datacenter) (err error) {
	query := make(map[string]interface{})
	query["names"] = ids
	if err := NewBaseModel("datacenter").FindBy(query, ds); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all entities on the store current user
// has access to
func (d *Datacenter) FindAll(au User, datacenters *[]Datacenter) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("datacenter").FindBy(query, datacenters); err != nil {
		return err
	}
	return nil
}

// Save : calls datacenter.set with the marshalled current entity
func (d *Datacenter) Save() (err error) {
	if err := NewBaseModel("datacenter").Save(d); err != nil {
		return err
	}

	return nil
}

// Delete : will delete a datacenter by its id
func (d *Datacenter) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = d.ID
	if err := NewBaseModel("datacenter").Delete(query); err != nil {
		return err
	}
	return nil
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (d *Datacenter) Redact() {
	d.AccessKeyID = ""
	d.SecretAccessKey = ""
	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	if d.Username != "" {
		d.Username, _ = crypto.Decrypt(d.Username, key)
	}
	d.Password = ""
}

// Improve : adds extra data to this entity
func (d *Datacenter) Improve() {
}

// Services : Get the services related with current datacenter
func (d *Datacenter) Services() (services []Service, err error) {
	var s Service
	err = s.FindByDatacenterID(d.ID, &services)

	return services, err
}

// GetID : ID getter
func (d *Datacenter) GetID() string {
	return d.Name
}

// GetType : Gets the resource type
func (d *Datacenter) GetType() string {
	return "project"
}

// Override : override not empty parameters with the given datacenter ones
func (d *Datacenter) Override(dt Datacenter) {
	if dt.Region != "" {
		d.Region = dt.Region
	}
	if dt.Username != "" {
		d.Username = dt.Username
	}
	if dt.Password != "" {
		d.Password = dt.Password
	}
	if dt.VCloudURL != "" {
		d.VCloudURL = dt.VCloudURL
	}
	if dt.VseURL != "" {
		d.VseURL = dt.VseURL
	}
	if dt.ExternalNetwork != "" {
		d.ExternalNetwork = dt.ExternalNetwork
	}
	if dt.AccessKeyID != "" {
		d.AccessKeyID = dt.AccessKeyID
	}
	if dt.SecretAccessKey != "" {
		d.SecretAccessKey = dt.SecretAccessKey
	}
	if dt.SubscriptionID != "" {
		d.SubscriptionID = dt.SubscriptionID
	}
	if dt.ClientID != "" {
		d.ClientID = dt.ClientID
	}
	if dt.ClientSecret != "" {
		d.ClientSecret = dt.ClientSecret
	}
	if dt.TenantID != "" {
		d.TenantID = dt.TenantID
	}
	if dt.Environment != "" {
		d.Environment = dt.Environment
	}
}
