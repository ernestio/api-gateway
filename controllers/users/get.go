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
	var r models.Role
	var roles []models.Role
	var proles []models.Role

	if !au.IsAdmin() {
		if au.Username != u {
			return 404, models.NewJSONError("User not found")
		}
	}

	if err := au.FindByUserName(u, &user); err != nil {
		return 404, models.NewJSONError("User not found")
	}

	if err := r.FindAllByUserAndResource(user.GetID(), "project", &proles); err == nil {
		user.ProjectMemberships = proles
	}

	if err := r.FindAllByUserAndResource(user.GetID(), "environment", &roles); err == nil {
		user.EnvMemberships = roles
	}

	user.Redact()

	body, err := json.Marshal(user)
	if err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
