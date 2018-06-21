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

	if existing.ResourceType == "project" {
		var roles []models.Role
		var owner bool

		err := existing.FindAllByResource(existing.ResourceID, existing.ResourceType, &roles)
		if err != nil {
			return 500, models.NewJSONError(err.Error())
		}

		for _, v := range roles {
			if v.Role == "owner" && v.UserID != existing.UserID {
				owner = true
			}
		}

		if !owner {
			return 400, models.NewJSONError("Cannot remove the only project owner")
		}
	}

	if err := existing.Delete(); err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, models.NewJSONError("Role deleted")
}
