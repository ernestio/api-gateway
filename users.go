package main

import (
	"encoding/json"
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
	GroupID  string `json:"group_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

// Validate the user
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Password == "" {
		return errors.New("User password is empty")
	}

	if u.GroupID == "" {
		return errors.New("User group is empty")
	}

	return nil
}

// Map : maps a user from a request's body and validates the input
func (u *User) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return ErrBadReqBody
	}

	err = u.Validate()
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

func getUsersHandler(c echo.Context) error {
	msg, err := n.Request("user.get", nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getUserHandler(c echo.Context) error {
	subject := fmt.Sprintf("user.get.%s", c.Param("user"))
	msg, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if len(msg.Data) == 0 {
		return ErrNotFound
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func createUserHandler(c echo.Context) error {
	var u User

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("user.create", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateUserHandler(c echo.Context) error {
	var u User

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	// Check if authenticated user is admin or updating itself
	au := authenticatedUser(c)
	if au.Username != u.Username && au.Admin != true {
		return ErrUnauthorized
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("user.update", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteUserHandler(c echo.Context) error {
	subject := fmt.Sprintf("user.delete.%s", c.Param("user"))
	_, err := n.Request(subject, nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.String(http.StatusOK, "")
}
