package users

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /users/:id: by deleting an
// existing user
func Delete(au models.User, user string) (int, []byte) {
	if err := au.Delete(user); err != nil {
		return 404, models.NewJSONError("User not found")
	}

	return http.StatusOK, []byte(`{"status": "User successfully deleted"}`)
}
