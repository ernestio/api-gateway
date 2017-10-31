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

	if !au.IsAdmin() {
		if au.Username != u {
			return 404, []byte("User not found")
		}
	}

	if err := au.FindByUserName(u, &user); err != nil {
		return 404, []byte("User not found")
	}

	if err := r.FindAllByUserAndResource(user.GetID(), "project", &roles); err == nil {
		for _, v := range roles {
			user.Projects = append(user.Projects, v.ResourceID+" ("+v.Role+")")
		}
	}
	if err := r.FindAllByUserAndResource(user.GetID(), "environment", &roles); err == nil {
		for _, v := range roles {
			user.Envs = append(user.Envs, v.ResourceID+" ("+v.Role+")")
		}
	}

	user.Redact()

	body, err := json.Marshal(user)
	if err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
