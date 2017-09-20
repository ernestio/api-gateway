package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ernestio/mapping/definition"
	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
)

// Given an echo context, it will extract the json or yml
// request body and will processes it in order to extract
// a valid defintion
func mapInputBuild(c echo.Context) (definition definition.Definition, raw []byte, err error) {
	req := c.Request()
	raw, err = ioutil.ReadAll(req.Body)

	// Normalize input body to json
	ctype := req.Header.Get("Content-Type")

	if ctype != "application/json" && ctype != "application/yaml" {
		return definition, raw, errors.New(`"Invalid input format"`)
	}

	if ctype == "application/yaml" {
		raw, err = yaml.YAMLToJSON(raw)
		if err != nil {
			return definition, raw, errors.New(`"Invalid yaml input"`)
		}
	}

	if err = json.Unmarshal(raw, &definition); err != nil {
		return definition, raw, errors.New(`"Invalid input"`)
	}

	// Override name and project if they're provided on the url
	if c.Param("project") != "" {
		definition["project"] = c.Param("project")
	}
	if c.Param("env") != "" {
		definition["name"] = c.Param("env")
	}

	return definition, raw, nil
}
