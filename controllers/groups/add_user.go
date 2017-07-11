package groups

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// AddUser : Adds an user to a group
func AddUser(au models.User, g string, body []byte) (int, []byte) {
	var err error
	var group models.Group
	var user models.User
	var payload map[string]string

	if err := group.FindByName(g, &group); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err := user.FindByUserName(payload["username"], &user); err != nil {
		return 400, []byte(err.Error())
	}

	user.GroupID = group.ID
	user.Password = ""
	user.Salt = ""
	if err := user.Save(); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("User " + user.Username + " successfully added to group " + group.Name)
}
