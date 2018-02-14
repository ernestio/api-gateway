package projects

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /projects/ by creating a
// project on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var err error
	var d models.Project
	var existing models.Project

	if d.Map(body) != nil {
		return 400, models.NewJSONError("Input is not valid")
	}

	err = d.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}

	if err := existing.FindByName(d.Name); err == nil {
		return 409, models.NewJSONError("Specified project already exists")
	}

	if err = d.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}
	if err := au.SetOwner(&d); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
