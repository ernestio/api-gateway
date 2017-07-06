package groups

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /groups/ by creating a group
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var g models.Group
	var existing models.Group
	var err error

	if g.Map(body) != nil {
		return 400, []byte("Input is not valid")
	}

	if err := existing.FindByName(g.Name, &existing); err == nil {
		return 409, []byte("Specified group already exists")
	}

	if err = g.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(g); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
