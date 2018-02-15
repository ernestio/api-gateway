package roles

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /roles/:id: by deleting an
// existing role
func Delete(au models.User, body []byte) (int, []byte) {
	var d models.Role

	if d.Map(body) != nil {
		return 400, models.NewJSONError("Input is not valid")
	}

	err := d.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, models.NewJSONError(err.Error())
	}
	if !au.IsAdmin() {
		if ok := au.IsOwner(d.ResourceType, d.ResourceID); !ok {
			return 403, models.NewJSONError("You're not authorized to perform this action")
		}
	}

	existing, err := d.Get(d.UserID, d.ResourceID, d.ResourceType)
	if !(err != nil || existing != nil) {
		return 409, models.NewJSONError("Specified role does not exists")
	}
	if err := existing.Delete(); err != nil {
		return 500, models.NewJSONError(err.Error())
	}

	return http.StatusOK, []byte(`{"status": "Role successfully deleted"}`)
}
