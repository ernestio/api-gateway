package roles

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /roles/ with a list of all
// roles
func List(au models.User) (int, []byte) {
	var err error
	var roles []models.Role
	var body []byte
	var r models.Role

	if err = r.FindAll(&roles); err != nil {
		return 404, models.NewJSONError(err.Error())
	}

	if body, err = json.Marshal(roles); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}
	return http.StatusOK, body
}
