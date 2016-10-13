/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/labstack/echo"
)

// Service holds the service response from service-store
type Service struct {
	ID           string      `json:"id"`
	GroupID      int         `json:"group_id"`
	DatacenterID int         `json:"datacenter_id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Version      time.Time   `json:"version"`
	Options      string      `json:"options"`
	Status       string      `json:"status"`
	Endpoint     string      `json:"endpoint"`
	Definition   interface{} `json:"definition"`
	Maped        string      `json:"mapping"`
}

type ServiceMapping struct {
	Vpcs struct {
		Items []struct {
			VpcID string `json:"vpc_id"`
		} `json:"items"`
	} `json:"vpcs"`
	Networks struct {
		Items []struct {
			Name   string `json:"name"`
			Subnet string `json:"network_aws_id"`
		} `json:"items"`
	} `json:"networks"`
	Instances struct {
		Items []struct {
			Name          string `json:"name"`
			InstanceAWSID string `json:"instance_aws_id"`
			PublicIP      string `json:"public_ip"`
			IP            string `json:"ip"`
		} `json:"items"`
	} `json:"instances"`
	Nats struct {
		Items []struct {
			Name            string `json:"name"`
			NatGatewayAWSID string `json:"nat_gateway_aws_id"`
		} `json:"items"`
	} `json:"nats"`
	SecurityGroups struct {
		Items []struct {
			Name               string `json:"name"`
			SecurityGroupAWSID string `json:"security_group_aws_id"`
		} `json:"items"`
	} `json:"firewalls"`
	Elbs struct {
		Items []struct {
			Name    string `json:"name"`
			DNSName string `json:"dns_name"`
		} `json:"items"`
	} `json:"elbs"`
}

// Validate the service
func (s *Service) Validate() error {
	if s.Name == "" {
		return errors.New("Service name is empty")
	}

	if s.DatacenterID == 0 {
		return errors.New("Service group is empty")
	}

	if s.Type == "" {
		return errors.New("Service type is empty")
	}

	return nil
}

// Map : maps a service from a request's body and validates the input
func (s *Service) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return ErrBadReqBody
	}

	err = s.Validate()
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

// Find : Searches for all services with filters
func (s *Service) Find(query map[string]interface{}, services *[]Service) (err error) {
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// FindByName : Searches for all services with a name equal to the specified
func (s *Service) FindByName(name string, service *Service) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	if err := NewBaseModel("service").GetBy(query, service); err != nil {
		return err
	}
	return nil
}

// FindByGroupID : Searches for all services on the store current user
// has access to with the specified group id
func (s *Service) FindByGroupID(id int, services *[]Service) (err error) {
	query := make(map[string]interface{})
	query["group_id"] = id

	return NewBaseModel("service").FindBy(query, services)
}

// FindByNameAndGroupID : Searches for all services with a name equal to the specified
func (s *Service) FindByNameAndGroupID(name string, id int, service *[]Service) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	query["group_id"] = id

	return NewBaseModel("service").FindBy(query, service)
}

// FindByID : Gets a model by its id
func (s *Service) FindByID(id int) (err error) {
	query := make(map[string]interface{})
	query["id"] = id
	if err := NewBaseModel("service").GetBy(query, s); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all groups on the store current user
// has access to
func (s *Service) FindAll(au User, services *[]Service) (err error) {
	query := make(map[string]interface{})
	query["group_id"] = au.GroupID
	if err := NewBaseModel("service").FindBy(query, services); err != nil {
		return err
	}
	return nil
}

// Save : calls service.set with the marshalled current group
func (s *Service) Save() (err error) {
	if err := NewBaseModel("service").Save(s); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a service by its id
func (s *Service) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = s.ID
	if err := NewBaseModel("service").Delete(query); err != nil {
		return err
	}
	return nil
}

// Mapping : will get a service mapping
func (s *Service) Mapping() (m ServiceMapping, err error) {
	query := make(map[string]interface{})
	query["id"] = s.ID

	err = NewBaseModel("service").callStoreBy("get.mapping", query, &m)

	return m, err
}

// Reset : will reset the service status to errored
func (s *Service) Reset() (err error) {
	s.Status = "errored"
	query := make(map[string]interface{})
	query["id"] = s.ID
	query["status"] = "errored"

	err = NewBaseModel("service").Set(query)

	return err
}
