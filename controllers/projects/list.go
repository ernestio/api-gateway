package projects

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /projects/ with a list of all
// projects
func List(au models.User) (int, []byte) {
	var err error
	var projects []models.Project
	var body []byte

	projects, err = au.GetProjects()

	if err != nil {
		return 404, []byte(err.Error())
	}

	for i := 0; i < len(projects); i++ {
		projects[i].Redact()
		projects[i].Improve()
	}

	if body, err = json.Marshal(projects); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
