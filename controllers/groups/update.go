package groups

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : responds to PUT /groups/:id: by updating an existing
// group
func Update(au models.User, body []byte) (int, []byte) {
	var g models.Group
	var existing models.Group
	var err error

	if g.Map(body) != nil {
		return 400, []byte("Invalid input")
	}

	if au.Admin != true {
		return 403, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action")
	}

	if err := existing.FindByName(g.Name, &existing); err != nil {
		return 404, []byte("Specified group does not exists")
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
