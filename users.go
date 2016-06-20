/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	"github.com/labstack/echo"
)

const (
	SaltSize = 32
	HashSize = 64
)

// User holds the user response from user-store
type User struct {
	ID       int    `json:"id"`
	GroupID  int    `json:"group_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Admin    bool   `json:"admin"`
}

// Validate vaildate all of the user's input
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("User username is empty")
	}

	if u.Password == "" {
		return errors.New("User password is empty")
	}

	if u.GroupID == 0 {
		return errors.New("User group is empty")
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

// HashPassword generates a secure salt and hashes the user's password
func (u *User) HashPassword() error {
	salt := make([]byte, SaltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	hash, err := scrypt.Key([]byte(u.Password), salt, 16384, 8, 1, HashSize)
	if err != nil {
		return err
	}

	// Create a base64 string of the binary salt and hash for storage
	u.Salt = base64.StdEncoding.EncodeToString(salt)
	u.Password = base64.StdEncoding.EncodeToString(hash)

	return nil
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
	msg, err := n.Request("user.find", nil, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusOK, msg.Data)
}

func getUserHandler(c echo.Context) error {
	query := fmt.Sprintf(`{"username": "%s"}`, c.Param("user"))
	msg, err := n.Request("user.get", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
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

	// Generate salt and hash password
	err := u.HashPassword()
	if err != nil {
		return ErrInternal
	}

	data, err := json.Marshal(u)
	if err != nil {
		return ErrInternal
	}

	msg, err := n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func updateUserHandler(c echo.Context) error {
	var u User
	var qu User

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
	msg, err := n.Request("user.get", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	err = json.Unmarshal(msg.Data, &qu)
	if err != nil {
		return ErrInternal
	}

	if qu.ID == 0 {
		return ErrNotFound
	}

	// update the user
	msg, err = n.Request("user.set", data, 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.JSONBlob(http.StatusAccepted, msg.Data)
}

func deleteUserHandler(c echo.Context) error {
	query := fmt.Sprintf(`{"username": "%s"}`, c.Param("user"))
	msg, err := n.Request("user.del", []byte(query), 5*time.Second)
	if err != nil {
		return ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return re.HTTPError
	}

	return c.String(http.StatusOK, "")
}
