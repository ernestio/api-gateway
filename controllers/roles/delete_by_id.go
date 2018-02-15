package roles

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// DeleteByID : responds to DELETE /roles/:id: by deleting an
// existing role
func DeleteByID(au models.User, id string) (int, []byte) {
	var err error
	var existing models.Role

	if err = existing.FindByID(id, &existing); err != nil {
		return 404, models.NewJSONError("Not found")
	}

	if err := existing.Delete(); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, models.NewJSONError("Role deleted")
}
