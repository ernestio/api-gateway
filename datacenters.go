package main

import (
	"errors"
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
	msg, err := n.Request("datacenter.find", nil, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getDatacenterHandler(c echo.Context) error {
	// TODO : Validate the datacenter is owned by this user
	body := []byte(`{"name":"` + c.Param("datacenter") + `"}`)
	msg, err := n.Request("datacenter.get", body, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if len(msg.Data) == 0 {
		return ErrNotFound
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createDatacenterHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	msg, err := n.Request("datacenter.set", data, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateDatacenterHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	msg, err := n.Request("datacenter.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteDatacenterHandler(c echo.Context) error {
	// TODO : Validate the datacenter is owned by this user
	body := []byte(`{"name":"` + c.Param("datacenter") + `"}`)
	_, err := n.Request("datacenter.del", body, 1*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.String(http.StatusOK, "")
}
