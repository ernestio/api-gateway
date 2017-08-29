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
		return 400, []byte(`{"code":400, "message":"` + err.Error() + `"}`)
	}

	u.Type = "local"

	if len(u.Password) < 8 {
		return 400, []byte(`Minimum password length is 8 characters`)
	}

	if err := existing.FindByUserName(u.Username, &existing); err == nil {
		return 409, []byte(`Specified user already exists`)
	}

	if err := u.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Error creating user")
	}

	u.Redact()

	body, err := json.Marshal(u)
	if err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
