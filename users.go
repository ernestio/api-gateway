/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	"github.com/labstack/echo"
)

const (
	// SaltSize is the lenght of the salt string
	SaltSize = 32
	// HashSize is the lenght of the hash string
	HashSize = 64
)

// User holds the user response from user-store
type User struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldpassword,omitempty"`
	Salt        string `json:"salt,omitempty"`
	Admin       bool   `json:"admin"`
}

// Validate vaildate all of the user's input
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Password == "" {
		return errors.New("User password is empty")
	}

	return nil
}

// Map a user from a request's body and validates the input
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

// Redact removes all sensitive fields from the return data before outputting to the user
func (u *User) Redact() {
	u.Password = ""
	u.Salt = ""
}

// ValidPassword checks if a submitted password matches the users password hash
func (u *User) ValidPassword(pw string) bool {
	userpass, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return false
	}

	usersalt, err := base64.StdEncoding.DecodeString(u.Salt)
	if err != nil {
		return false
	}

	hash, err := scrypt.Key([]byte(pw), usersalt, 16384, 8, 1, HashSize)
	if err != nil {
		return false
	}

	// Compare in constant time to mitigate timing attacks
	if subtle.ConstantTimeCompare(userpass, hash) == 1 {
		return true
	}

	return false
}

func getUsersHandler(c echo.Context) error {
	var query string
	var users []User

	au := authenticatedUser(c)

	if !au.Admin {
		query = fmt.Sprintf(`{"group_id": %d}`, au.GroupID)
	}

	msg, err := n.Request("user.find", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	// Remove sensitive data
	err = json.Unmarshal(msg.Data, &users)
	if err != nil {
		return ErrInternal
	}

	for i := 0; i < len(users); i++ {
		users[i].Redact()
	}

	return c.JSON(http.StatusOK, users)
}

func getUserHandler(c echo.Context) error {
	var query string
	var user User

	au := authenticatedUser(c)

	if au.Admin {
		query = fmt.Sprintf(`{"id": %s}`, c.Param("user"))
	} else {
		query = fmt.Sprintf(`{"id": %s, "group_id": %d}`, c.Param("user"), au.GroupID)
	}

	msg, err := n.Request("user.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	// Remove sensitive data
	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		return ErrInternal
	}

	user.Redact()

	return c.JSON(http.StatusOK, user)
}

func createUserHandler(c echo.Context) error {
	var u User
	var created User
	var existing User

	if authenticatedUser(c).Admin != true {
		return ErrUnauthorized
	}

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	// Check if the user exists
	query := fmt.Sprintf(`{"username": "%s"}`, u.Username)
	msg, err := n.Request("user.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	err = json.Unmarshal(msg.Data, &existing)
	if err != nil {
		return ErrInternal
	}

	if existing.ID != 0 {
		c.Response().Header().Add("Location", fmt.Sprintf("/users/%d", existing.ID))
		return ErrExists
	}

	// Create the user
	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	msg, err = n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	// Remove sensitive data
	err = json.Unmarshal(msg.Data, &created)
	if err != nil {
		return ErrInternal
	}

	created.Redact()

	return c.JSON(http.StatusOK, created)
}

func updateUserHandler(c echo.Context) error {
	var u User
	var updated User
	var existing User

	if u.Map(c) != nil {
		return ErrBadReqBody
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	// Check if authenticated user is admin or updating itself
	au := authenticatedUser(c)
	if au.Username != u.Username && au.Admin != true {
		return ErrUnauthorized
	}

	// Check user exists
	query := fmt.Sprintf(`{"id": %d}`, u.ID)
	msg, err := n.Request("user.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	err = json.Unmarshal(msg.Data, &existing)
	if err != nil {
		return ErrInternal
	}

	if existing.ID == 0 {
		return ErrNotFound
	}

	// Check a non-admin user is not trying to change their group
	if au.Admin != true && u.GroupID != existing.GroupID {
		return ErrUnauthorized
	}

	// Check the old password if it is present
	if u.OldPassword != "" && !existing.ValidPassword(u.OldPassword) {
		return ErrUnauthorized
	}

	// update the user
	msg, err = n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	// Remove sensitive data
	err = json.Unmarshal(msg.Data, &updated)
	if err != nil {
		return ErrInternal
	}

	updated.Redact()

	return c.JSON(http.StatusOK, updated)
}

func deleteUserHandler(c echo.Context) error {
	au := authenticatedUser(c)
	if au.Admin != true {
		return ErrUnauthorized
	}

	query := fmt.Sprintf(`{"id": %s}`, c.Param("user"))
	msg, err := n.Request("user.del", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		fmt.Println(re.Error)
		return re.HTTPError
	}

	return c.String(http.StatusOK, "")
}
