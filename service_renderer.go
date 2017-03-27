/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"

	graph "gopkg.in/r3labs/graph.v2"
)

// ServiceRender : Service representation to be rendered on the frontend
type ServiceRender struct {
	ID             string              `json:"id"`
	DatacenterID   int                 `json:"datacenter_id"`
	Name           string              `json:"name"`
	Version        string              `json:"version"`
	Status         string              `json:"status"`
	UserID         int                 `json:"user_id"`
	UserName       string              `json:"user_name"`
	LastKnownError string              `json:"last_known_error"`
	Options        string              `json:"options"`
	Definition     string              `json:"definition"`
	Vpcs           []map[string]string `json:"vpcs"`
	Networks       []map[string]string `json:"networks"`
	Instances      []map[string]string `json:"instances"`
	Nats           []map[string]string `json:"nats"`
	SecurityGroups []map[string]string `json:"security_groups"`
	Elbs           []map[string]string `json:"elbs"`
	RDSClusters    []map[string]string `json:"rds_clusters"`
	RDSInstances   []map[string]string `json:"rds_instances"`
	EBSVolumes     []map[string]string `json:"ebs_volumes"`
}

// Render : Map a Service to a ServiceRender
func (o *ServiceRender) Render(s Service) error {
	o.ID = s.ID
	o.DatacenterID = s.DatacenterID
	o.Name = s.Name
	o.Version = s.Version.String()
	o.Status = s.Status
	o.UserID = s.UserID
	o.UserName = s.UserName
	if def, ok := s.Definition.(string); ok == true {
		o.Definition = def
	}

	g, err := s.Mapping()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	o.Vpcs = RenderVpcs(g)
	o.Networks = RenderNetworks(g)
	o.SecurityGroups = RenderSecurityGroups(g)
	o.Nats = RenderNats(g)
	o.Instances = RenderInstances(g)
	o.Elbs = RenderELBs(g)
	o.RDSClusters = RenderRDSClusters(g)
	o.RDSInstances = RenderRDSInstances(g)
	o.EBSVolumes = RenderEBSVolumes(g)

	return err
}

// RenderVpcs : renders a services vpcs
func RenderVpcs(g *graph.Graph) []map[string]string {
	var vpcs []map[string]string

	for _, n := range g.GetComponents().ByType("vpc") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["vpc_aws_id"].(string)
		subnet, _ := (*gc)["vpc_subnet"].(string)
		vpcs = append(vpcs, map[string]string{
			"name":       name,
			"vpc_id":     id,
			"vpc_subnet": subnet,
		})
	}

	return vpcs
}

// RenderNetworks : renders a services networks
func RenderNetworks(g *graph.Graph) []map[string]string {
	var networks []map[string]string

	for _, n := range g.GetComponents().ByType("network") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["network_aws_id"].(string)
		az, _ := (*gc)["availablity_zone"].(string)
		networks = append(networks, map[string]string{
			"name":              name,
			"network_aws_id":    id,
			"availability_zone": az,
		})
	}

	return networks
}

// RenderSecurityGroups : renders a services security groups
func RenderSecurityGroups(g *graph.Graph) []map[string]string {
	var sgs []map[string]string

	for _, n := range g.GetComponents().ByType("firewall") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["security_group_aws_id"].(string)
		sgs = append(sgs, map[string]string{
			"name":                  name,
			"security_group_aws_id": id,
		})
	}

	return sgs
}

// RenderNats : renders a services nat gateways
func RenderNats(g *graph.Graph) []map[string]string {
	var nats []map[string]string

	for _, n := range g.GetComponents().ByType("nat_gateway") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["nat_gateway_aws_id"].(string)
		nats = append(nats, map[string]string{
			"name":               name,
			"nat_gateway_aws_id": id,
		})
	}

	return nats
}

// RenderELBs : renders a services elbs
func RenderELBs(g *graph.Graph) []map[string]string {
	var elbs []map[string]string

	for _, n := range g.GetComponents().ByType("elb") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		dns, _ := (*gc)["dns_name"].(string)
		elbs = append(elbs, map[string]string{
			"name":     name,
			"dns_name": dns,
		})
	}

	return elbs
}

// RenderInstances : renders a services instances
func RenderInstances(g *graph.Graph) []map[string]string {
	var instances []map[string]string

	for _, n := range g.GetComponents().ByType("instance") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["instance_aws_id"].(string)
		pip, _ := (*gc)["public_ip"].(string)
		ip, _ := (*gc)["ip"].(string)
		instances = append(instances, map[string]string{
			"name":            name,
			"instance_aws_id": id,
			"public_ip":       pip,
			"ip":              ip,
		})
	}

	return instances
}

// RenderRDSClusters : renders a services rds clusters
func RenderRDSClusters(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("rds_cluster") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		endpoint, _ := (*gc)["endpoint"].(string)
		rdss = append(rdss, map[string]string{
			"name":     name,
			"endpoint": endpoint,
		})
	}

	return rdss
}

// RenderRDSInstances : renders a services rds instances
func RenderRDSInstances(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("rds_instance") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		endpoint, _ := (*gc)["endpoint"].(string)
		rdss = append(rdss, map[string]string{
			"name":     name,
			"endpoint": endpoint,
		})
	}

	return rdss
}

// RenderEBSVolumes : renders a services ebs volumes
func RenderEBSVolumes(g *graph.Graph) []map[string]string {
	var rdss []map[string]string

	for _, n := range g.GetComponents().ByType("ebs_volume") {
		gc := n.(*graph.GenericComponent)
		name, _ := (*gc)["name"].(string)
		id, _ := (*gc)["volume_aws_id"].(string)
		rdss = append(rdss, map[string]string{
			"name":          name,
			"volume_aws_id": id,
		})
	}

	return rdss
}

// RenderCollection : Maps a collection of Service on a collection of ServiceRender
func (o *ServiceRender) RenderCollection(services []Service) (list []ServiceRender, err error) {
	for _, s := range services {
		var output ServiceRender
		if err := output.Render(s); err == nil {
			list = append(list, output)
		}
	}

	return list, nil
}

// ToJSON : Converts a ServiceRender to json string
func (o *ServiceRender) ToJSON() ([]byte, error) {
	return json.Marshal(o)
}
