package datacenters

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /datacenters/ by creating a
// datacenter on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var err error
	var d models.Datacenter
	var existing models.Datacenter

	if au.GroupID == 0 {
		return 401, []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action")
	}

	if d.Map(body) != nil {
		return 400, []byte("Input is not valid")
	}

	err = d.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	d.GroupID = au.GroupID

	if err := existing.FindByName(d.Name, &existing); err == nil {
		return 409, []byte("Specified datacenter already exists")
	}

	if err = d.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
