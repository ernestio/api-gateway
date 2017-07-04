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
	var group models.Group
	var err error

	if au.Admin == true {
		if err := group.FindAll(au, &groups); err != nil {
			h.L.Warning(err.Error())
		}
	} else {
		if err := group.FindByID(au.GroupID); err != nil {
			h.L.Warning(err.Error())
		}
		groups = append(groups, group)
	}

	if body, err = json.Marshal(groups); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
