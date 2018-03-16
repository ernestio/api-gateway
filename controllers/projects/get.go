package projects

import (
	"encoding/json"
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /projects/:id:/ with the specified
// project details
func Get(au models.User, project string) (int, []byte) {
	var envs []models.Env
	var r models.Role
	var roles []models.Role
	var d models.Project
	var body []byte
	var err error

	if st, res := h.IsAuthorizedToResource(&au, h.GetProject, d.GetType(), project); st != 200 {
		return st, res
	}

	if err := d.FindByName(project); err != nil {
		return 404, models.NewJSONError("Project not found")
	}

	query := make(map[string]interface{}, 0)
	query["datacenter_id"] = d.ID
	envs, err = au.EnvsBy(query)
	if err == nil {
		for _, v := range envs {
			nameParts := strings.Split(v.Name, models.EnvNameSeparator)
			if nameParts[0] == project {
				d.Environments = append(d.Environments, nameParts[1])
			}
		}
	}

	if err := r.FindAllByResource(d.GetID(), d.GetType(), &roles); err == nil {
		for _, v := range roles {
			d.Roles = append(d.Roles, v.UserID+" ("+v.Role+")")
		}
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
