package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Datacenter holds the datacenter response from datacenter-store
type Datacenter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate the datacenter
func (d *Datacenter) Validate() error {
	if d.Name == "" {
		return errors.New("Datacenter name is empty")
	}

	return nil
}

func getDatacentersHandler(c echo.Context) error {
	msg, err := n.Request("datacenters.get", nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getDatacenterHandler(c echo.Context) error {
	subject := fmt.Sprintf("datacenters.get.%s", c.Param("datacenter"))
	msg, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	if len(msg.Data) == 0 {
		return notFound
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createDatacenterHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("datacenters.create", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateDatacenterHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("datacenters.update", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteDatacenterHandler(c echo.Context) error {
	subject := fmt.Sprintf("datacenters.delete.%s", c.Param("datacenter"))
	_, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.String(http.StatusOK, "")
}
