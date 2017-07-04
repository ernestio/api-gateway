package services

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Search : Finds all services
func Search(au models.User, query map[string]interface{}) (int, []byte) {
	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	b, err := json.Marshal(list)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	return http.StatusOK, b
}
