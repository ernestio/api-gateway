package envs

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

// CreateBuild : Creates a build for a specific environment
func CreateBuild(au models.User, s models.ServiceInput, definition, body []byte, dry string) (int, []byte) {
	var err error
	var group []byte
	var previous models.Env
	var mapping map[string]interface{}
	var prevID string

	dt := models.Project{
		Credentials: make(map[string]interface{}),
	}

	// *********** VALIDATIONS *********** //

	if parts := strings.Split(s.Name, models.EnvNameSeparator); len(parts) > 2 {
		return 400, []byte("Environment name does not support char '" + models.EnvNameSeparator + "' as part of its name")
	}

	// Get datacenter
	if err = dt.FindByName(s.Datacenter, &dt); err != nil {
		h.L.Error(err.Error())
		return 400, []byte("Specified project does not exist")
	}

	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get previous env if exists
	previous, _ = previous.FindLastByName(s.Name)
	if &previous != nil {
		prevID = previous.ID
		if previous.Status == "in_progress" {
			h.L.Error("Environment is still in progress")
			return http.StatusNotFound, []byte(`"Your environment process is 'in progress' if your're sure you want to fix it please reset it first"`)
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
	pcredentials := models.Project{
		Credentials: make(map[string]interface{}),
	}

	if &previous != nil {
		if previous.Credentials != nil {
			prevDT := models.Project{
				Credentials: previous.Credentials,
			}
			pcredentials.Override(prevDT)
		}
	}

	if s.Credentials != nil {
		newDT := models.Project{
			Credentials: s.Credentials,
		}

		newDT.Encrypt()
		pcredentials.Override(newDT)
	}

	dt.Override(pcredentials)
	rawDatacenter, err := json.Marshal(dt)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error trying to get the project")
	}

	// *********** REQUESTING DEFINITION ************ //

	payload := ServicePayload{
		ID:         generateEnvID(s.Name + "-" + s.Datacenter),
		PrevID:     prevID,
		Service:    (*json.RawMessage)(&body),
		Datacenter: (*json.RawMessage)(&rawDatacenter),
		Group:      (*json.RawMessage)(&group),
	}

	if body, err = json.Marshal(payload); err != nil {
		return 500, []byte("Internal server error")
	}
	var def models.Definition
	mapping, err = def.MapCreation(body)

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

	// *********** SAVE NEW ENV AND PROCESS CREATION *********** //
	ss := models.Env{
		ID:           payload.ID,
		Name:         s.Name,
		Type:         dt.Type,
		UserID:       currentUser.ID,
		DatacenterID: dt.ID,
		Version:      time.Now(),
		Status:       "in_progress",
		Definition:   d,
		Mapped:       mapping,
		Credentials:  pcredentials.Credentials,
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
	err = ss.RequestCreation(mapping)

	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	parts := strings.Split(s.Name, "/")

	return http.StatusOK, []byte(`{"id":"` + payload.ID + `", "project": "` + parts[0] + `",  "name":"` + parts[1] + `"}`)
}
