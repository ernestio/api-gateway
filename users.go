package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// User holds the user response from user-store
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// Validate the user
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("User name is empty")
	}

	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Username == "" {
		return errors.New("User password is empty")
	}

	return nil
}

func getUsersHandler(c echo.Context) error {
	msg, err := n.Request("users.get", nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getUserHandler(c echo.Context) error {
	subject := fmt.Sprintf("users.get.%s", c.Param("user"))
	msg, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	if len(msg.Data) == 0 {
		return notFound
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createUserHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("users.create", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateUserHandler(c echo.Context) error {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return badReqBody
	}

	msg, err := n.Request("users.update", data, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteUserHandler(c echo.Context) error {
	subject := fmt.Sprintf("users.delete.%s", c.Param("user"))
	_, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return gatewayTimeout
	}

	return c.String(http.StatusOK, "")
}
