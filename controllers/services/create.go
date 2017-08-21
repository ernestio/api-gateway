package services

import (
	"encoding/json"
	"net/http"
	"strings"
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

// Create : Will receive a service application
func Create(au models.User, s models.ServiceInput, definition, body []byte, isAnImport bool, dry string) (int, []byte) {
	var err error
	var group []byte
	var previous models.Service
	var mapping map[string]interface{}
	var prevID string
	var dt models.Datacenter

	// *********** VALIDATIONS *********** //

	if parts := strings.Split(s.Name, models.EnvNameSeparator); len(parts) > 2 {
		return 400, []byte("Environment name does not support char '" + models.EnvNameSeparator + "' as part of its name")
	}

	// Get datacenter
	if err = dt.FindByName(s.Datacenter, &dt); err != nil {
		h.L.Error(err.Error())
		return 400, []byte("Specified datacenter does not exist")
	}

	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get previous service if exists
	previous, _ = previous.FindLastByName(s.Name)
	if &previous != nil {
		prevID = previous.ID
		if previous.Status == "in_progress" {
			h.L.Error("Service is still in progress")
			return http.StatusNotFound, []byte(`"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
		}
	}
	if prevID == "" {
		if st, res := h.IsAuthorizedToResource(&au, h.UpdateProject, dt.GetType(), s.Datacenter); st != 200 {
			return st, res
		}
	} else {
		if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, previous.GetType(), s.Name); st != 200 {
			return st, res
		}
	}

	// *********** OVERRIDE PROJECT CREDENTIALS ************ //
	credentials := models.Datacenter{}
	if &previous != nil {
		if previous.ProjectInfo != nil {
			var prevDT models.Datacenter
			if err := json.Unmarshal(*previous.ProjectInfo, &prevDT); err == nil {
				credentials.Override(prevDT)
			}
		}
	}

	if s.ProjectInfo != nil {
		var newDT models.Datacenter
		if err := json.Unmarshal(*s.ProjectInfo, &newDT); err == nil {
			newDT.Encrypt()
			credentials.Override(newDT)
		}
	}

	dt.Override(credentials)
	rawDatacenter, err := json.Marshal(dt)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error trying to get the datacenter")
	}
	rawCredentials, err := json.Marshal(credentials)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error trying to get the datacenter")
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
		mapping, err = def.MapImport(body)
	} else {
		mapping, err = def.MapCreation(body)
	}

	if err != nil {
		h.L.Error(err.Error())
		return 400, []byte(err.Error())
	}

	// *********** BUILD REQUEST IF IS DRY *********** //

	if dry == "true" {
		res, err := views.RenderDefinition(mapping)
		if err != nil {
			h.L.Error(err.Error())
			return 400, []byte("Internal error")
		}
		return http.StatusOK, res
	}

	d := string(definition)
	if defParts := strings.Split(d, "credentials:"); len(defParts) > 0 {
		d = defParts[0]
	}

	// *********** SAVE NEW SERVICE AND PROCESS CREATION / IMPORT *********** //
	ss := models.Service{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         dt.Type,
		UserID:       currentUser.ID,
		DatacenterID: dt.ID,
		Version:      time.Now(),
		Status:       "in_progress",
		Definition:   d,
		Mapped:       mapping,
		ProjectInfo:  (*json.RawMessage)(&rawCredentials),
	}

	if err := ss.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	if prevID == "" {
		if err := au.SetOwner(&ss); err != nil {
			return 500, []byte("Internal server error")
		}
	}

	// Apply changes
	if isAnImport == true {
		err = ss.RequestImport(mapping)
	} else {
		err = ss.RequestCreation(mapping)
	}

	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + payload.ID + `", "name":"` + s.Name + `"}`)
}
