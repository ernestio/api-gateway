package projects

import (
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /projects/:id: by deleting an
// existing project
func Delete(au models.User, project string) (int, []byte) {
	var d models.Project
	var err error

	if err = d.FindByName(project); err != nil {
		id, err := strconv.Atoi(project)
		if err = d.FindByID(id); err != nil {
			return 404, []byte("Project not found")
		}
	}

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteProject, d.GetType(), d.Name); st != 200 {
		return st, res
	}

	ss, err := d.Envs()
	if err != nil {
		return 500, []byte(err.Error())
	}

	if len(ss) > 0 {
		return 400, []byte("Existing environments are referring to this project.")
	}

	if err := d.Delete(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte("Project successfully deleted")
}
