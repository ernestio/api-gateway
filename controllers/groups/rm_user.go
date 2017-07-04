package groups

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// RmUser : Deletes an user from a group
func RmUser(au models.User, u string) (int, []byte) {
	var user models.User

	if au.Admin == false {
		return http.StatusForbidden, []byte("You don't have permissions to perform this action, please login with an admin account")
	}

	if err := user.FindByID(u, &user); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}
	user.GroupID = 0
	user.Password = ""
	user.Salt = ""
	if err := user.Save(); err != nil {
		return http.StatusGatewayTimeout, []byte("Internal server error")
	}

	return http.StatusOK, []byte("User " + user.Username + " successfully removed from group")
}
