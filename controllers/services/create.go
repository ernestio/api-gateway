package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// ServicePayload : payload to be sent to workflow manager
type ServicePayload struct {
	ID         string           `json:"id"`
	PrevID     string           `json:"previous_id"`
	Datacenter *json.RawMessage `json:"datacenter"`
	Group      *json.RawMessage `json:"client"`
	Service    *json.RawMessage `json:"service"`
}

// CreateServiceHandler : Will receive a service application
func CreateServiceHandler(au models.User, s models.ServiceInput, definition, body []byte, isAnImport bool, dry string) (int, []byte) {
	var err error
	var group []byte
	var previous *models.Service
	var service []byte
	var prevID string
	var dt models.Datacenter

	// *********** VALIDATIONS *********** //

	// Get datacenter
	dt, err = dt.GetByNameAndGroupID(s.Datacenter, au.GroupID)
	if err != nil {
		h.L.Error(err.Error())
		return 400, []byte(err.Error())
	}

	rawDatacenter, err := json.Marshal(dt)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error trying to get the datacenter")
	}

	// Get group
	if group, err = getRawGroup(au.GroupID); err != nil {
		h.L.Error(err.Error())
		return http.StatusNotFound, []byte(err.Error())
	}

	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get previous service if exists
	if previous, err = previous.GetByNameAndGroupID(s.Name, au.GroupID); err != nil {
		h.L.Error("Previous service not found")
		return http.StatusNotFound, []byte(err.Error())
	}

	if previous != nil {
		prevID = previous.ID
		if previous.Status == "in_progress" {
			h.L.Error("Service is still in progress")
			return http.StatusNotFound, []byte(`"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
		}
	}

	// *********** REQUESTING DEFINITION ************ //

	payload := ServicePayload{
		ID:         generateServiceID(s.Name + "-" + s.Datacenter),
		PrevID:     prevID,
		Service:    (*json.RawMessage)(&body),
		Datacenter: (*json.RawMessage)(&rawDatacenter),
		Group:      (*json.RawMessage)(&group),
	}

	if body, err = json.Marshal(payload); err != nil {
		return 500, []byte("Internal server error")
	}
	var def models.Definition
	if isAnImport == true {
		service, err = def.MapImport(body)
	} else {
		service, err = def.MapCreation(body)
	}

	if err != nil {
		h.L.Error(err.Error())
		return 400, []byte(err.Error())
	}

	// *********** BUILD REQUEST IF IS DRY *********** //

	if dry == "true" {
		res, err := views.RenderDefinition(service)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("Internal error")
		}
		return http.StatusOK, res
	}

	// *********** SAVE NEW SERVICE AND PROCESS CREATION / IMPORT *********** //

	ss := models.Service{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         dt.Type,
		GroupID:      au.GroupID,
		UserID:       currentUser.ID,
		DatacenterID: dt.ID,
		Version:      time.Now(),
		Status:       "in_progress",
		Definition:   string(definition),
		Maped:        string(service),
	}

	if err := ss.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	// Apply changes
	if isAnImport == true {
		err = ss.RequestImport(service)
	} else {
		err = ss.RequestCreation(service)
	}

	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + payload.ID + `", "name":"` + s.Name + `"}`)
}

func getRawGroup(id int) (group []byte, err error) {
	var g models.Group

	if err = g.FindByID(id); err != nil {
		return group, errors.New(`"Specified group does not exist"`)
	}

	if group, err = json.Marshal(g); err != nil {
		return group, errors.New(`"Internal error"`)
	}
	h.L.Info(group)

	return group, nil
}
