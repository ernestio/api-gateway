package users

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /users/ by creating a user
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var u models.User
	var existing models.User

	if err := u.Map(body); err != nil {
		h.L.Error(err.Error())
		return 400, models.NewJSONError(`{"code":400, "message":"` + err.Error() + `"}`)
	}

	err := u.Validate()
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	if err := existing.FindByUserName(u.Username, &existing); err == nil {
		return 409, models.NewJSONError(`Specified user already exists`)
	}

	u.Type = "local"

	if err := u.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Error creating user")
	}

	if u.MFA != nil {
		if *u.MFA {
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
