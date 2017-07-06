package services

import (
	"encoding/json"
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
	var datacenter []byte
	var err error
	var group []byte
	var previous *models.Service

	payload := ServicePayload{}

	if au.GroupID == 0 {
		body := "Current user does not belong to any group."
		body += "\nPlease assign the user to a group before performing this action"
		return 401, []byte(body)
	}

	payload.Service = (*json.RawMessage)(&body)

	// Get datacenter
	if datacenter, err = getDatacenter(s.Datacenter, au.GroupID); err != nil {
		h.L.Error(err.Error())
		return 404, []byte(err.Error())
	}
	payload.Datacenter = (*json.RawMessage)(&datacenter)

	// Get group
	if group, err = getGroup(au.GroupID); err != nil {
		h.L.Error(err.Error())
		return http.StatusNotFound, []byte(err.Error())
	}
	payload.Group = (*json.RawMessage)(&group)
	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Generate service ID
	payload.ID = generateServiceID(s.Name + "-" + s.Datacenter)

	// Get previous service if exists
	if previous, err = getService(s.Name, au.GroupID); err != nil {
		h.L.Error("Previous service not found")
		return http.StatusNotFound, []byte(err.Error())
	}

	if previous != nil {
		payload.PrevID = previous.ID
		if previous.Status == "in_progress" {
			h.L.Error("Service is still in progress")
			return http.StatusNotFound, []byte(`"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
		}
	}

	var service []byte

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

	if dry == "true" {
		res, err := views.RenderDefinition(service)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("Internal error")
		}
		return http.StatusOK, res
	}

	var datacenterStruct struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(datacenter, &datacenterStruct); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	ss := models.Service{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         datacenterStruct.Type,
		GroupID:      au.GroupID,
		UserID:       currentUser.ID,
		DatacenterID: datacenterStruct.ID,
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

	return http.StatusOK, []byte(`{"id":"` + payload.ID + `"}`)
}
