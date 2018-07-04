package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	err := u.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}

	uid, _ := strconv.Atoi(user)

	if u.Username != user && u.ID != uid {
		return 400, models.NewJSONError("User does not match payload name")
	}

	// Check if authenticated user is admin or updating itself
	if !u.CanBeChangedBy(au) {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, models.NewJSONError(err.Error())
	}

	// Check user exists
	if err := au.FindByUserName(user, &existing); err != nil {
		if err := au.FindByID(user, &existing); err != nil {
			h.L.Error(err.Error())
			return 404, models.NewJSONError("Specified user not found")
		}
	}

	if existing.ID == 0 {
		err := errors.New("Specified user not found")
		h.L.Error(err.Error())
		return 404, models.NewJSONError(err.Error())
	}

	u.Username = existing.Username

	if !au.IsAdmin() && existing.Username != au.Username {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, models.NewJSONError(err.Error())
	}

	if !au.IsAdmin() && existing.IsAdmin() != u.IsAdmin() {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, models.NewJSONError(err.Error())
	}

	if u.Password != nil {
		err := u.Validate()
		if err != nil {
			return 400, models.NewJSONError(err.Error())
		}

		// Check the old password if it is present
		if u.OldPassword != nil && !existing.ValidPassword(*u.OldPassword) {
			err := errors.New("Provided credentials are not valid")
			h.L.Error(err.Error())
			return 403, models.NewJSONError(err.Error())
		}
	}

	if err := u.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Error updating user")
	}

	if existing.MFA == nil {
		existing.MFA = h.Bool(false)
	}

	if u.MFA != nil {
		if *u.MFA && !*existing.MFA {
			mfaSecret := u.MFASecret
			u.Redact(au)
			u.MFASecret = mfaSecret
		} else {
			u.Redact(au)
		}
	} else {
		u.Redact(au)
	}

	body, err = json.Marshal(u)
	if err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
