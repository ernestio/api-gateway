package groups

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /groups/ with a list of all
// groups
func List(au models.User) (int, []byte) {
	var groups []models.Group
	var body []byte
	var err error

	if groups, err = au.Groups(); err != nil {
		h.L.Warning(err.Error())
		return 400, []byte(err.Error())
	}

	if body, err = json.Marshal(groups); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
