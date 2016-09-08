package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
	"github.com/nats-io/nats"
	"github.com/nu7hatch/gouuid"
)

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter string `json:"datacenter"`
	Name       string `json:"name"`
}

type ServicePayload struct {
	ID         string           `json:"id"`
	PrevID     string           `json:"previous_id"`
	Datacenter *json.RawMessage `json:"datacenter"`
	Group      *json.RawMessage `json:"client"`
	Service    *json.RawMessage `json:"service"`
}

// Maps input as a valid Serviceinput
func mapInputService(c echo.Context) (s ServiceInput, definition []byte, jsonbody []byte, err error) {
	req := c.Request()
	definition, err = ioutil.ReadAll(req.Body())

	// Normalize input body to json
	ctype := req.Header().Get("Content-Type")

	if ctype != "application/json" && ctype != "application/yaml" {
		return s, definition, jsonbody, errors.New(`"Invalid input format"`)
	}

	if ctype == "application/yaml" {
		jsonbody, err = yaml.YAMLToJSON(definition)
		if err != nil {
			return s, definition, jsonbody, errors.New(`"Invalid yaml input"`)
		}
	} else {
		jsonbody = definition
	}

	if err = json.Unmarshal(jsonbody, &s); err != nil {
		return s, definition, jsonbody, errors.New(`"Invalid input"`)
	}

	return s, definition, jsonbody, nil
}

// Generates a service ID based on an input service
func generateServiceID(salt string) string {
	sufix := generateStreamID(salt)
	prefix, _ := uuid.NewV4()

	return prefix.String() + "-" + string(sufix[:])
}

func generateStreamID(salt string) string {
	compose := []byte(salt)
	hasher := md5.New()
	hasher.Write(compose)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getDatacenter(name string, group int) (datacenter []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"name": "%s", "group_id": %d}`, name, group)
	if msg, err = n.Request("datacenter.find", []byte(query), 1*time.Second); err != nil {
		return datacenter, ErrGatewayTimeout
	}
	if string(msg.Data) == `[]` {
		return datacenter, errors.New(`"Specified datacenter does not exist"`)
	}

	// Get only the first datcenter
	datacenters := make([]interface{}, 0)
	json.Unmarshal(msg.Data, &datacenters)
	res, err := json.Marshal(datacenters[0])
	if err != nil {
		return datacenter, errors.New("Internal error trying to get the datacenter")
	}

	return res, nil
}

func getDatacenterByID(id int) (user []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"id": %d}`, id)
	if msg, err = n.Request("datacenter.get", []byte(query), 1*time.Second); err != nil {
		return user, ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return user, errors.New(`"Specified datacenter does not exist"`)
	}
	return msg.Data, nil
}

func getUser(id int) (user []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"id": %d}`, id)
	if msg, err = n.Request("user.get", []byte(query), 1*time.Second); err != nil {
		return user, ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return user, errors.New(`"Specified user does not exist"`)
	}
	return msg.Data, nil
}

func getGroup(id int) (group []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"id": %d}`, id)
	if msg, err = n.Request("group.get", []byte(query), 1*time.Second); err != nil {
		return group, ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return group, errors.New(`"Specified group does not exist"`)
	}
	return msg.Data, nil
}

func getService(name string, group int) (service *Service, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"name":"%s","group_id":%d}`, name, group)
	if msg, err = n.Request("service.find", []byte(query), 1*time.Second); err != nil {
		return service, ErrGatewayTimeout
	}

	p := []Service{}
	json.Unmarshal(msg.Data, &p)
	if len(p) == 0 {
		return nil, nil
	}

	return &p[0], nil
}

func mapCreateDefinition(payload ServicePayload) (body []byte, err error) {
	var msg *nats.Msg

	if body, err = json.Marshal(payload); err != nil {
		return body, errors.New("Provided yaml is not valid")
	}

	if msg, err = n.Request("definition.map.creation", body, 1*time.Second); err != nil {
		return body, errors.New("Provided yaml is not valid")
	}

	var s struct {
		Error string `json:"error"`
	}

	json.Unmarshal(msg.Data, &s)
	if s.Error != "" {
		return body, errors.New(s.Error)
	}

	return msg.Data, nil
}

func getServiceRaw(name string, group int) (service []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"name":"%s","group_id":%d}`, name, group)
	if msg, err = n.Request("service.find", []byte(query), 1*time.Second); err != nil {
		return service, ErrGatewayTimeout
	}
	p := []*json.RawMessage{}

	if err = json.Unmarshal(msg.Data, &p); err != nil {
		return nil, errors.New(`"Internal error"`)
	}

	if len(p) == 0 {
		return nil, errors.New(`"Service not found"`)
	}

	if body, err := p[0].MarshalJSON(); err != nil {
		return nil, errors.New("Internal error")
	} else {
		return body, nil
	}
}

type OutputService struct {
	ID           string `json:"id"`
	DatacenterID int    `json:"datacenter_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
	Options      string `json:"options"`
	Endpoint     string `json:"endpoint"`
	Definition   string `json:"definition"`
	Networks     []struct {
		Name   string `json:"name"`
		Subnet string `json:"network_aws_id"`
	} `json:"networks"`
	Instances []struct {
		Name          string `json:"name"`
		InstanceAWSID string `json:"instance_aws_id"`
		PublicIP      string `json:"public_ip"`
		IP            string `json:"ip"`
	} `json:"instances"`
	SecurityGroups []struct {
		Name            string `json:"name"`
		NatGatewayAWSID string `json:"nat_gateway_aws_id"`
	} `json:"nats"`
}

type ServiceMapping struct {
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
	SecurityGroups struct {
		Items []struct {
			Name            string `json:"name"`
			NatGatewayAWSID string `json:"nat_gateway_aws_id"`
		} `json:"items"`
	} `json:"nats"`
}

func getServicesOutput(filter map[string]interface{}) (list []OutputService, err error) {
	var msg *nats.Msg

	query, err := json.Marshal(filter)
	if err != nil {
		return list, err
	}

	if msg, err = n.Request("service.find", query, 1*time.Second); err != nil {
		return list, ErrGatewayTimeout
	}

	if err := json.Unmarshal(msg.Data, &list); err != nil {
		return list, errors.New("Internal error")
	}

	// Popolate service with detailedd info
	for i, v := range list {
		mapping := ServiceMapping{}
		if msg, err = n.Request("service.get.mapping", []byte(`{"id":"`+v.ID+`"}`), 1*time.Second); err != nil {
			return list, ErrGatewayTimeout
		}

		if err := json.Unmarshal(msg.Data, &mapping); err != nil {
			return list, errors.New("Internal error")
		}

		list[i].Networks = mapping.Networks.Items
		list[i].SecurityGroups = mapping.SecurityGroups.Items
		list[i].Instances = mapping.Instances.Items
	}

	return list, nil
}

func resetService(au User, name string) (status int, err error) {
	var list []OutputService
	filter := make(map[string]interface{})
	filter["group_id"] = au.GroupID
	filter["name"] = name

	if list, err = getServicesOutput(filter); err != nil {
		return 500, errors.New("Internal error")
	}
	if len(list) == 0 {
		return 404, errors.New(`No services found with for '` + name + `'`)
	}
	if list[0].Status != "in_progress" {
		return 200, errors.New("Reset only applies to 'in progress' serices, however service '" + name + "' is on status '" + list[0].Status)
	}

	query := `{"id":"` + list[0].ID + `","status":"errored"}`
	if _, err := n.Request("service.set", []byte(query), 1*time.Second); err != nil {
		return 500, errors.New("Could not update the service")
	}

	return 200, nil
}
func saveService(id string, name string, t string, v time.Time, s string, o string, d string, m string, group uint, datacenter uint) {
	var payload struct {
		Uuid         string    `json:"id"`
		GroupID      uint      `json:"group_id"`
		DatacenterID uint      `json:"datacenter_id"`
		Name         string    `json:"name"`
		Type         string    `json:"type"`
		Version      time.Time `json:"version"`
		Status       string    `json:"status"`
		Options      string    `json:"options"`
		Definition   string    `json:"definition"`
		Mapping      string    `json:"mapping"`
	}

	payload.Uuid = id
	payload.Name = name
	payload.Type = t
	payload.GroupID = group
	payload.DatacenterID = datacenter
	payload.Version = v
	payload.Status = s
	payload.Options = o
	payload.Definition = d
	payload.Mapping = m

	body, _ := json.Marshal(payload)

	n.Publish("service.set", body)
}
