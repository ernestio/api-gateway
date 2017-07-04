package users

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /users/:id:/ with the specified
// user details
func Get(au models.User, u string) (int, []byte) {
	var user models.User

	if err := au.FindByID(u, &user); err != nil {
		return 404, []byte("User not found")
	}
	user.Redact()

	body, err := json.Marshal(user)
	if err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
