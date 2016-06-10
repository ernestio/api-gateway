package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func getUsersHander(c echo.Context) error {
	msg, err := n.Request("users.get", nil, 5*time.Second)
	if err != nil {
		c.Error(err)
	}
	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getUserHander(c echo.Context) error {
	subject := fmt.Sprintf("users.get.%s", c.Param("user"))
	msg, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		c.Error(err)
	}
	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createUserHander(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		c.Error(err)
	}
	msg, err := n.Request("users.create", data, 5*time.Second)
	if err != nil {
		c.Error(err)
	}
	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateUserHander(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		c.Error(err)
	}
	msg, err := n.Request("users.update", data, 5*time.Second)
	if err != nil {
		c.Error(err)
	}
	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteUserHander(c echo.Context) error {
	subject := fmt.Sprintf("users.delete.%s", c.Param("user"))
	_, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		c.Error(err)
	}
	return c.String(http.StatusOK, `{"status": "ok"}`)
}
