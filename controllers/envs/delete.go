package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : Deletes a service by name
func Delete(au models.User, name string) (int, []byte) {
	var err error
	var def models.Definition
	var s models.Env
	var dt models.Project

	if s, err = s.FindLastByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if s.ID == "" {
		println("innnn")
		return 404, []byte("Specified environment name does not exist")
	}

	// Get datacenter
	if err = dt.FindByID(s.DatacenterID); err != nil {
		h.L.Error(err.Error())
		return 400, []byte("Specified project does not exist")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteEnv, s.GetType(), name); st != 200 {
		return st, res
	}

	if s.Status == "in_progress" {
		return 400, []byte(`"Environment is already applying some changes, please wait until they are done"`)
	}

	credentials := models.Project{}
	if s.ProjectInfo != nil {
		var newDT models.Project
		if err := json.Unmarshal(*s.ProjectInfo, &newDT); err == nil {
			credentials.Override(newDT)
		}
	}

	dt.Override(credentials)
	rawDatacenter, err := json.Marshal(dt)
	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error trying to get the project")
	}

	query := []byte(`{"previous_id":"` + s.ID + `","datacenter":` + string(rawDatacenter) + `}`)
	//++++++++++++++++++
	body, err := def.MapDeletion(query)

	if err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't map the environment"`)
	}
	if err := s.RequestDeletion(body); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(`"Couldn't call service.delete"`)
	}

	return http.StatusOK, []byte(`{"id":"` + s.ID + `"}`)
}
