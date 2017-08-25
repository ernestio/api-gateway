package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ernestio/api-gateway/models"
	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
)

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

	// Override name and project if they're provided on the url
	if c.Param("project") != "" {
		s.Datacenter = c.Param("project")
	}
	if c.Param("env") != "" {
		s.Name = c.Param("env")
	}

	return s, definition, jsonbody, nil
}
