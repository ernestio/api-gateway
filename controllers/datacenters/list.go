package datacenters

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /datacenters/ with a list of all
// datacenters
func List(au models.User) (int, []byte) {
	var err error
	var datacenters []models.Datacenter
	var body []byte
	var datacenter models.Datacenter

	if au.Admin == true {
		err = datacenter.FindAll(au, &datacenters)
	} else {
		datacenters, err = au.Datacenters()
	}

	if err != nil {
		return 404, []byte(err.Error())
	}

	for i := 0; i < len(datacenters); i++ {
		datacenters[i].Redact()
		datacenters[i].Improve()
	}

	if body, err = json.Marshal(datacenters); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
