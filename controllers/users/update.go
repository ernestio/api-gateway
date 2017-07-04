package users

import (
	"encoding/json"
	"errors"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : responds to PUT /users/:id: by updating an existing
// user
func Update(au models.User, user string, body []byte) (int, []byte) {
	var u models.User
	var existing models.User

	if err := u.Map(body); err != nil {
		h.L.Error(err.Error())
		return 400, []byte(err.Error())
	}

	if len(u.Password) < 8 {
		err := errors.New("Minimum password length is 8 characters")
		h.L.Error(err.Error())
		return 400, []byte(`{"code":400, "message":"` + err.Error() + `"}`)
	}

	// Check if authenticated user is admin or updating itself
	if au.Username != u.Username && au.Admin != true {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	// Check user exists
	if err := au.FindByID(user, &existing); err != nil {
		h.L.Error(err.Error())
		return 404, []byte("Specified user not found")
	}

	if existing.ID == 0 {
		err := errors.New("Specified user not found")
		h.L.Error(err.Error())
		return 404, []byte(err.Error())
	}

	// Check a non-admin user is not trying to change their group
	if au.Admin != true && u.GroupID != existing.GroupID {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	// Check the old password if it is present
	if u.OldPassword != "" && !existing.ValidPassword(u.OldPassword) {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	if err := u.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Error updating user")
	}

	u.Redact()

	body, err := json.Marshal(u)
	if err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
