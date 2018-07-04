/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
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
	ID                 int     `json:"id"`
	Username           string  `json:"username"`
	Password           *string `json:"password,omitempty"`
	OldPassword        *string `json:"oldpassword,omitempty"`
	Salt               string  `json:"salt,omitempty"`
	Admin              *bool   `json:"admin,omitempty"`
	MFA                *bool   `json:"mfa,omitempty"`
	MFASecret          string  `json:"mfa_secret,omitempty"`
	VerificationCode   string  `json:"verification_code,omitempty"`
	EnvMemberships     []Role  `json:"env_memberships,omitempty"`
	ProjectMemberships []Role  `json:"project_memberships,omitempty"`
	Type               string  `json:"type,omitempty"`
	Disabled           *bool   `json:"disabled,omitempty"`
}

// AuthResponse : Describes an Authenticator service response
type AuthResponse struct {
	OK      bool   `json:"ok"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

// Authenticate verifies user credentials
func (u *User) Authenticate() (*AuthResponse, error) {
	if !IsAlphaNumeric(u.Username) {
		return nil, errors.New("bad username specified")
	}

	mfa, err := u.IsMFA()
	if err != nil {
		return nil, err
	}

	if mfa && u.VerificationCode == "" {
		return nil, errors.New("mfa required")
	}

	msg, err := N.Request("authentication.get", []byte(`{"username": "`+u.Username+`", "password": "`+*u.Password+`", "verification_code": "`+u.VerificationCode+`"}`), 10*time.Second)
	if err != nil {
		return nil, err
	}

	res := AuthResponse{}
	err = json.Unmarshal(msg.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Validate checks user input details for missing values and invalid characters
func (u *User) Validate() error {
	r := regexp.MustCompile(`^[a-zA-Z0-9@._\-]*$`)
	if u.Username == "" {
		return errors.New("Username cannot be empty")
	}
	if !r.MatchString(u.Username) {
		return errors.New(`Username can only contain the following characters: a-z 0-9 @._-`)
	}
	if u.Password != nil {
		if *u.Password == "" {
			return errors.New("Password cannot be empty")
		}
		if len(*u.Password) < 8 {
			return errors.New("Minimum password length is 8 characters")
		}
		if !r.MatchString(*u.Password) {
			return errors.New(`Password can only contain the following characters: a-z 0-9 @._-`)
		}
	}
	return nil
}

// Map a user from a request's body and validates the input
func (u *User) Map(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}
	return nil
}

// FindByUserName : find a user for the given username, and maps it on
// the fiven User struct
func (u *User) FindByUserName(name string, user *User) (err error) {
	query := make(map[string]interface{})
	query["username"] = name
	return NewBaseModel("user").GetBy(query, user)
}

// FindAll : Searches for all users on the store current user
// has access to
func (u *User) FindAll(users *[]User) (err error) {
	query := make(map[string]interface{})
	if !u.IsAdmin() {
		// TODO add auth
	}
	return NewBaseModel("user").FindBy(query, users)
}

// FindByID : Searches a user by ID on the store current user
// has access to
func (u *User) FindByID(id string, user *User) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	if !u.IsAdmin() {
		// TODO add auth
	}
	return NewBaseModel("user").GetBy(query, user)
}

// Save : calls user.set with the marshalled current user
func (u *User) Save() (err error) {
	return NewBaseModel("user").Save(u)
}

// Delete : will delete a user by its id
func (u *User) Delete(id string) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	return NewBaseModel("user").Delete(query)
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (u *User) Redact(au User) {
	empty := ""
	u.Password = &empty
	u.Salt = ""
	u.MFASecret = ""

	if !au.IsAdmin() {
		u.Admin = nil
		u.MFA = nil
		u.Password = nil
		u.Disabled = nil
		u.Type = ""
	}
}

// Improve : adds extra data
func (u *User) Improve() {
}

// ValidPassword : checks if a submitted password matches
// the users password hash
func (u *User) ValidPassword(pw string) bool {
	userpass, err := base64.StdEncoding.DecodeString(*u.Password)
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

// GetPolicies : Gets the related user policies if any
func (u *User) GetPolicies() (ds []Policy, err error) {
	var d Policy

	if u.IsAdmin() {
		err = d.FindAll(&ds)
	} else {
		var r Role
		if names, err := r.FindAllIDsByUserAndType(u.GetID(), d.GetType()); err == nil {
			if names == nil {
				return ds, nil
			}
			err = d.FindByNames(names, &ds)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}

	return ds, err
}

// GetProjects : Gets the related user projects if any
func (u *User) GetProjects() (ds []Project, err error) {
	var d Project

	if u.IsAdmin() {
		err = d.FindAll(*u, &ds)
	} else {
		var r Role
		if ids, err := r.FindAllIDsByUserAndType(u.GetID(), d.GetType()); err == nil {
			if ids == nil {
				return ds, nil
			}
			err = d.FindByIDs(ids, &ds)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}

	return ds, err
}

// ProjectByName : Gets the related user projects if any
func (u *User) ProjectByName(name string) (d Project, err error) {
	if err = d.FindByName(name); err != nil {
		err = errors.New("Project not found")
	}

	return
}

// FindAllKeyValue : Finds all users on a id:name hash
func (u *User) FindAllKeyValue() (list map[int]string) {
	var users []User
	list = make(map[int]string)
	if err := u.FindAll(&users); err != nil {
		h.L.Warning(err.Error())
	}
	for _, v := range users {
		list[v.ID] = v.Username
	}
	return list
}

// GetBuild : Gets a specific build if authorized
func (u *User) GetBuild(id string) (build Env, err error) {
	var envs []Env
	var s Env

	query := make(map[string]interface{})
	query["id"] = id
	err = s.Find(query, &envs)

	if len(envs) == 0 {
		h.L.Debug("Build " + id + " not found")
		return build, errors.New("Not found")
	}

	return
}

// EnvsBy : Get authorized envs by any filter
func (u *User) EnvsBy(filters map[string]interface{}) ([]Env, error) {
	var uEnvs []Env
	var err error
	var e Env
	var rEnvs []Env

	if u.IsAdmin() {
		err = e.Find(filters, &uEnvs)
		if err != nil {
			return nil, err
		}

		return uEnvs, err
	} else {
		var envs []Env
		var r Role
		var roles []Role

		err = e.Find(nil, &envs)
		if err != nil {
			return nil, err
		}

		err = r.FindAllByUser(u.Username, &roles)
		if err != nil {
			return nil, nil
		}

		for _, r := range roles {
			if r.ResourceType != "project" && r.ResourceType != "environment" {
				continue
			}

			for _, e := range envs {
				if r.ResourceType == "project" {
					if r.ResourceID == strings.Split(e.Name, "/")[0] {
						rEnvs = append(rEnvs, e)
					}
				} else {
					if r.ResourceID == e.Name {
						rEnvs = append(rEnvs, e)
					}
				}
			}
		}
	}

	m := make(map[string]bool)

	for _, v := range rEnvs {
		if _, ok := m[v.Name]; !ok {
			uEnvs = append(uEnvs, v)
			m[v.Name] = true
		}
	}

	return uEnvs, nil
}

// CanBeChangedBy : Checks if an user has write permissions on another user
func (u *User) CanBeChangedBy(user User) bool {
	if user.IsAdmin() {
		return true
	}

	if u.Username == user.Username {
		return true
	}

	return false
}

// GetAdmin : admin getter
func (u *User) GetAdmin() bool {
	return u.IsAdmin()
}

// GetID : ID getter
func (u *User) GetID() string {
	return u.Username
}

type resource interface {
	GetID() string
	GetType() string
}

// SetOwner : ...
func (u *User) SetOwner(o resource) error {
	return u.setRole(o, "owner")
}

// SetReader : ...
func (u *User) SetReader(o resource) error {
	return u.setRole(o, "reader")
}

// setRole : ...
func (u *User) setRole(o resource, r string) error {
	role := Role{
		UserID:       u.GetID(),
		ResourceID:   o.GetID(),
		ResourceType: o.GetType(),
		Role:         r,
	}

	return role.Save()
}

// Owns : Checks if the user owns a specific resource
func (u *User) Owns(o resource) bool {
	return u.IsOwner(o.GetType(), o.GetID())

}

// IsOwner : check if is the owner of a specific resource
func (u *User) IsOwner(resourceType, resourceID string) bool {
	owner := "owner"
	reader := "reader"

	if u.IsAdmin() {
		return true
	}

	if role, err := u.getRole(resourceType, resourceID); err == nil {
		if role == owner {
			return true
		}
		if role == reader {
			return false
		}
	}
	if resourceType == "build" || resourceType == "environment" {
		parts := strings.Split(resourceID, "/")
		if len(parts) > 0 {
			if role, err := u.getRole("project", parts[0]); err == nil {
				if role == owner {
					return true
				}
			}
		}
	}

	return false
}

// IsReader : check if has reader permissions on a specific resource
func (u *User) IsReader(resourceType, resourceID string) bool {
	owner := "owner"
	reader := "reader"

	if u.IsAdmin() {
		return true
	}

	if role, err := u.getRole(resourceType, resourceID); err == nil {
		if role == reader || role == owner {
			return true
		}
	}
	if resourceType == "build" || resourceType == "environment" {
		parts := strings.Split(resourceID, "/")
		if len(parts) > 0 {
			if role, err := u.getRole("project", parts[0]); err == nil {
				if role == reader || role == owner {
					return true
				}
			}
		}
	}

	return false
}

func (u *User) getRole(resourceType, resourceID string) (string, error) {
	var role Role

	existing, err := role.Get(u.GetID(), resourceID, resourceType)
	if err != nil || existing == nil {
		return "", errors.New("Not found")
	}

	return existing.Role, nil
}

// IsAdmin : Check if a user is admin or not
func (u *User) IsAdmin() bool {
	if u.Admin != nil {
		return *u.Admin
	}
	return false
}

// IsMFA checks if a user has Multi-Factor authentication enabled
func (u *User) IsMFA() (bool, error) {
	msg, err := N.Request("user.get", []byte(`{"username": "`+u.Username+`"}`), time.Second)
	if err != nil {
		return false, err
	}

	var fetchedUser User
	err = json.Unmarshal(msg.Data, &fetchedUser)
	if err != nil {
		return false, err
	}

	if fetchedUser.MFA != nil {
		return *fetchedUser.MFA, nil
	}

	return false, nil
}
