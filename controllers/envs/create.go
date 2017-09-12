package envs

import (
	"net/http"
	"strings"
	"time"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : Creates an environment
func Create(au models.User, name, project string, credentials map[string]interface{}) (int, []byte) {
	var p models.Project

	if parts := strings.Split(name, models.EnvNameSeparator); len(parts) > 2 {
		return 400, []byte("Environment name does not support char '" + models.EnvNameSeparator + "' as part of its name")
	}

	envName := project + models.EnvNameSeparator + name
	// Get existing environment
	if _, err := getEnvRaw(au, envName); err == nil {
		return 404, []byte("Environment " + name + " already exists")
	}

	// Get datacenter
	if err := p.FindByName(project, &p); err != nil {
		h.L.Error(err.Error())
		return 400, []byte("Specified project does not exist")
	}

	var currentUser models.User
	if err := currentUser.FindByUserName(au.Username, &currentUser); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateProject, p.GetType(), project); st != 200 {
		return st, res
	}

	pcredentials := models.Project{
		Credentials: make(map[string]interface{}),
	}

	if credentials != nil {
		newDT := models.Project{
			Credentials: credentials,
		}

		newDT.Encrypt()
		pcredentials.Override(newDT)
	}

	env := models.Env{
		ID:           generateEnvID(envName),
		Name:         envName,
		Type:         p.Type,
		DatacenterID: p.ID,
		UserID:       currentUser.ID,
		Version:      time.Now(),
		Status:       "done",
		Credentials:  pcredentials.Credentials,
	}

	if err := env.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"id":"` + env.ID + `"}`)
}
