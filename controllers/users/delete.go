package users

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /users/:id: by deleting an
// existing user
func Delete(au models.User, user string) (int, []byte) {
	if au.Admin != true {
		return 403, []byte("You're not allowed to perform this action, please contact your admin")
	}

	if err := au.Delete(user); err != nil {
		return 404, []byte("User not found")
	}

	return http.StatusOK, []byte("User successfully deleted")
}
