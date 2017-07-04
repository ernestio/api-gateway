package services

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Builds : gets the list of builds for the specified service
func Builds(au models.User, query map[string]interface{}) (int, []byte) {
	var user models.User

	users := user.FindAllKeyValue()

	if au.Admin != true {
		query["group_id"] = au.GroupID
	}

	list, err := getServicesOutput(query)
	if err != nil {
		return 500, []byte(err.Error())
	}
	for i := range list {
		for id, name := range users {
			if id == list[i].UserID {
				list[i].UserName = name
			}
		}
	}

	body, err := json.Marshal(list)
	if err != nil {
		return 500, []byte("Internal error")
	}

	return http.StatusOK, body
}
