package projects

import (
	"encoding/json"
	"net/http"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : responds to PUT /projects/:id: by updating
// an existing project
func Update(au models.User, project string, body []byte) (int, []byte) {
	var d models.Project
	var existing models.Project
	var err error

	if d.Map(body) != nil {
		return 400, models.NewJSONError("Invalid input")
	}

	if err = existing.FindByName(project); err != nil {
		id, err := strconv.Atoi(project)
		if err = existing.FindByID(id); err != nil {
			return 404, models.NewJSONError("Project not found")
		}
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateProject, d.GetType(), d.Name); st != 200 {
		return st, res
	}

	existing.Credentials = d.Credentials

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
