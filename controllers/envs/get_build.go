package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// GetBuild : gets the details of a specific service build
func GetBuild(au models.User, query map[string]interface{}) (int, []byte) {
	var o views.ServiceRender

	builds, err := au.EnvsBy(query)
	if err != nil {
		return 500, []byte(err.Error())
	}

	list, err := o.RenderCollection(builds)
	if err != nil {
		return 500, []byte(err.Error())
	}

	if len(list) > 0 {
		if st, res := h.IsAuthorizedToResource(&au, h.GetBuild, builds[0].GetType(), builds[0].Name); st != 200 {
			return st, res
		}
		body, err := json.Marshal(list[0])
		if err != nil {
			return 500, []byte("Internal server error")
		}
		return http.StatusOK, body
	}
	return http.StatusNotFound, []byte("")
}
