package groups

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ernestio/api-gateway/models"
)

// Get : responds to GET /groups/:id:/ with the specified
// group details
func Get(group string) (int, []byte) {
	var err error
	var g models.Group
	var body []byte

	id, _ := strconv.Atoi(group)
	if err := g.FindByID(id); err != nil {
		return 404, []byte("Group not found")
	}

	if body, err = json.Marshal(g); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
