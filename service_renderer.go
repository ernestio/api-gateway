/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

type ServiceRender struct {
	ID             string `json:"id"`
	DatacenterID   int    `json:"datacenter_id"`
	Name           string `json:"name"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	UserID         int    `json:"user_id"`
	UserName       string `json:"user_name"`
	LastKnownError string `json:"last_known_error"`
	Options        string `json:"options"`
	Endpoint       string `json:"endpoint"`
	Definition     string `json:"definition"`
	VpcID          string `json:"vpc_id"`
	Networks       []struct {
		Name             string `json:"name"`
		Subnet           string `json:"network_aws_id"`
		AvailabilityZone string `json:"availability_zone"`
	} `json:"networks"`
	Instances []struct {
		Name          string `json:"name"`
		InstanceAWSID string `json:"instance_aws_id"`
		PublicIP      string `json:"public_ip"`
		IP            string `json:"ip"`
	} `json:"instances"`
	Nats []struct {
		Name            string `json:"name"`
		NatGatewayAWSID string `json:"nat_gateway_aws_id"`
	} `json:"nats"`
	SecurityGroups []struct {
		Name               string `json:"name"`
		SecurityGroupAWSID string `json:"security_group_aws_id"`
	} `json:"security_groups"`
	Elbs []struct {
		Name    string `json:"name"`
		DNSName string `json:"dns_name"`
	}
	RDSClusters []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	} `json:"rds_clusters"`
	RDSInstances []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	} `json:"rds_instances"`
}

func (o *ServiceRender) Render(s Service) (err error) {
	var mapping ServiceMapping

	o.ID = s.ID
	o.DatacenterID = s.DatacenterID
	o.Name = s.Name
	o.Version = s.Version.String()
	o.Status = s.Status
	o.UserID = s.UserID
	o.UserName = s.UserName
	o.Endpoint = s.Endpoint
	if def, ok := s.Definition.(string); ok == true {
		o.Definition = def
	}

	if mapping, err = s.Mapping(); err != nil {
		log.Println(err.Error())
		return err
	}
	if len(mapping.Vpcs.Items) > 0 {
		o.VpcID = mapping.Vpcs.Items[0].VpcID
	}

	o.LastKnownError = mapping.LastKnownError
	o.Networks = mapping.Networks.Items
	o.SecurityGroups = mapping.SecurityGroups.Items
	o.Nats = mapping.Nats.Items
	o.Instances = mapping.Instances.Items
	o.Elbs = mapping.Elbs.Items
	o.RDSClusters = mapping.RDSClusters.Items
	o.RDSInstances = mapping.RDSInstances.Items

	return err
}

func (o *ServiceRender) RenderCollection(services []Service) (list []ServiceRender, err error) {
	for _, s := range services {
		var output ServiceRender
		if err := output.Render(s); err == nil {
			list = append(list, output)
		}
	}

	return list, nil
}

func (o *ServiceRender) ToJson() ([]byte, error) {
	return json.Marshal(o)
}
