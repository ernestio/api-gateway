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
	"strings"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/scrypt"
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

// FindByUserName : find a user for the given username, and maps it on
// the fiven User struct
func (u *User) FindByUserName(name string, user *User) (err error) {
	var res []byte

	query := `{"username": "` + name + `"}`
	if res, err = Query("user.get", query); err != nil {
		return err
	}
	if strings.Contains(string(res), `"error"`) {
		return errors.New(`"Specified username does not exist"`)
	}
	if err = json.Unmarshal(res, &user); err != nil {
		return errors.New(`"Specified username does not exist"`)
	}

	return nil
}

// FindAll : Searches for all users on the store current user
// has access to
func (u *User) FindAll(users *[]User) (err error) {
	var query string
	var res []byte

	if !u.Admin {
		query = fmt.Sprintf(`{"group_id": %d}`, u.GroupID)
	}

	if res, err = Query("user.find", query); err != nil {
		return err
	}

	err = json.Unmarshal(res, &users)
	if err != nil {
		return ErrInternal
	}

	return nil
}

// FindByID : Searches a user by ID on the store current user
// has access to
func (u *User) FindByID(id string, user *User) (err error) {
	var query string
	var res []byte

	if u.Admin {
		query = fmt.Sprintf(`{"id": %s}`, id)
	} else {
		query = fmt.Sprintf(`{"id": %s, "group_id": %d}`, id, u.GroupID)
	}

	if res, err = Query("user.get", query); err != nil {
		return err
	}

	if err = json.Unmarshal(res, &user); err != nil {
		return ErrInternal
	}

	return nil
}

// Save : calls user.set with the marshalled current user
func (u *User) Save() (err error) {
	var res []byte

	data, err := json.Marshal(u)
	if err != nil {
		return ErrBadReqBody
	}

	if res, err = Query("user.set", string(data)); err != nil {
		return err
	}

	if err := json.Unmarshal(res, &u); err != nil {
		return ErrInternal
	}

	return nil
}

// Delete : will delete a user by its id
func (u *User) Delete(id string) (err error) {
	query := fmt.Sprintf(`{"id": %s}`, id)
	if _, err := Query("user.del", query); err != nil {
		return err
	}

	return nil
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (u *User) Redact() {
	u.Password = ""
	u.Salt = ""
}

// ValidPassword : checks if a submitted password matches
// the users password hash
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

// Group : Gets the related user group if any
func (u *User) Group() (group Group) {
	group.FindByID(u.GroupID)

	return group
}

// TODO : Move this to somewhere else
func Query(subject, query string) ([]byte, error) {
	var res []byte
	msg, err := n.Request(subject, []byte(query), 5*time.Second)
	if err != nil {
		return res, ErrGatewayTimeout
	}

	if re := responseErr(msg); re != nil {
		return res, re.HTTPError
	}

	return msg.Data, nil
}
