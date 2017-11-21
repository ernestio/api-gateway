package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ernestio/api-gateway/models"
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
	if err != nil {
		return definition, raw, err
	}

	// Normalize input body to json
	ctype := req.Header.Get("Content-Type")

	if ctype != "application/yaml" {
		return definition, raw, errors.New(`"Invalid input format"`)
	}

	if err = yaml.Unmarshal(raw, &definition); err != nil {
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

func mapAction(c echo.Context) (*models.Action, error) {
	var action models.Action

	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	return &action, json.Unmarshal(data, &action)
}
