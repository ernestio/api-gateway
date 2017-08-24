package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Builds : gets the list of builds for the specified service
func Builds(au models.User, name string) (int, []byte) {
	var o views.ServiceRender
	var s models.Env
	var builds []models.Env
	var err error

	if st, res := h.IsAuthorizedToResource(&au, h.ListBuilds, s.GetType(), name); st != 200 {
		return st, res
	}

	query := make(map[string]interface{}, 0)
	query["name"] = name
	if err = s.Find(query, &builds); err != nil {
		h.L.Warning(err.Error())
	}
	if err != nil {
		return 500, []byte(err.Error())
	}

	list, err := o.RenderCollection(builds)
	if err != nil {
		return 500, []byte(err.Error())
	}

	users := au.FindAllKeyValue()
	for i := range list {
		for id, name := range users {
			if id == list[i].UserID {
				list[i].UserName = name
			}
		}
	}

	body, err := json.Marshal(list)
	if err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}