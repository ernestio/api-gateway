package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Group holds the group response from group-store
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate the group
func (g *Group) Validate() error {
	if g.Name == "" {
		return errors.New("Group name is empty")
	}

	return nil
}

func getGroupsHandler(c echo.Context) error {
	msg, err := n.Request("groups.get", nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getGroupHandler(c echo.Context) error {
	subject := fmt.Sprintf("groups.get.%s", c.Param("group"))
	msg, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	if len(msg.Data) == 0 {
		return notFound
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createGroupHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("groups.create", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateGroupHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("groups.update", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteGroupHandler(c echo.Context) error {
	subject := fmt.Sprintf("groups.delete.%s", c.Param("group"))
	_, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.String(http.StatusOK, "")
}
