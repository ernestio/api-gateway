package services

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// GetBuild : gets the details of a specific service build
func GetBuild(au models.User, query map[string]interface{}) (int, []byte) {
	var err error
	var list []views.ServiceRender

	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	if list, err = getServicesOutput(query); err != nil {
		return 500, []byte(err.Error())
	}

	if len(list) > 0 {
		body, err := json.Marshal(list[0])
		if err != nil {
			return 500, []byte("Internal server error")
		}
		return http.StatusOK, body
	}
	return http.StatusNotFound, []byte("")
}
