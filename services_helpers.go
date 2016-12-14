package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

// ServicePayload : payload to be sent to workflow manager
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
	definition, err = ioutil.ReadAll(req.Body)

	// Normalize input body to json
	ctype := req.Header.Get("Content-Type")

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
	if _, err := hasher.Write(compose); err != nil {
		log.Println(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func getDatacenter(name string, group int) (datacenter []byte, err error) {
	var d Datacenter
	var datacenters []Datacenter

	if err := d.FindByNameAndGroupID(name, group, &datacenters); err != nil {
		return datacenter, err
	}

	if len(datacenters) == 0 {
		return datacenter, errors.New(`"Specified datacenter does not exist"`)
	}

	datacenter, err = json.Marshal(datacenters[0])
	if err != nil {
		return datacenter, errors.New("Internal error trying to get the datacenter")
	}

	return datacenter, nil
}

func getGroup(id int) (group []byte, err error) {
	var g Group

	if err = g.FindByID(id); err != nil {
		return group, errors.New(`"Specified group does not exist"`)
	}

	if group, err = json.Marshal(g); err != nil {
		return group, errors.New(`"Internal error"`)
	}
	println(group)

	return group, nil
}

func getService(name string, group int) (service *Service, err error) {
	var s Service
	var services []Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return service, ErrGatewayTimeout
	}

	if len(services) == 0 {
		return nil, nil
	}

	return &services[0], nil
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

	if err := json.Unmarshal(msg.Data, &s); err != nil {
		log.Println(err)
		return body, err
	}
	if s.Error != "" {
		return body, errors.New(s.Error)
	}

	return msg.Data, nil
}

func getServiceRaw(name string, group int) (service []byte, err error) {
	var s Service
	var services []Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return nil, errors.New(`"Internal error"`)
	}

	if len(services) == 0 {
		return nil, errors.New(`"Service not found"`)
	}

	body, err := json.Marshal(services[0])
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

func getServicesOutput(filter map[string]interface{}) (list []ServiceRender, err error) {
	var s Service
	var services []Service
	var o ServiceRender

	if err := s.Find(filter, &services); err != nil {
		return list, err
	}

	return o.RenderCollection(services)
}
