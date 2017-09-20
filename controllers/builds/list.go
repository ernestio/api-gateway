package builds

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /services/ with a list of all
// services for current user group
func List(au models.User, env string) (int, []byte) {
	var b models.Build
	var list []models.Build
	var body []byte

	err := b.FindByEnvironmentName(env, &list)
	if err != nil {
		h.L.Warning(err.Error())
		return 404, []byte("Environment not found")
	}

	body, err = json.Marshal(list)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
