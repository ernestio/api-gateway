package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ernestio/api-gateway/models"
	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
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

// Given an echo context, it will extract the json or yml
// request body and will processes it in order to extract
// a valid defintion
func mapInputService(c echo.Context) (s models.ServiceInput, definition []byte, jsonbody []byte, err error) {
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
