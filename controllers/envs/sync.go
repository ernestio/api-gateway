package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Sync : Respons to POST /services/:service/sync/ and synchronizes a service with
// its provider representation
func Sync(au models.User, name string) (int, []byte) {
	var raw []byte
	var err error
	var s models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.SyncEnv, s.GetType(), name); st != 200 {
		return st, res
	}

	// Get existing env
	if raw, err = getEnvRaw(au, name); err != nil {
		return 404, []byte(err.Error())
	}

	if err := json.Unmarshal(raw, &s); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	if s.Status == "in_progress" {
		return 400, []byte(`"Environment is already applying some changes, please wait until they are done"`)
	}

	if err = s.RequestSync(); err != nil {
		return 500, []byte("An error ocurred while ernest was trying to sync your environment")
	}

	// TODO : This probably needs to use the monit tool instead of this.

	return http.StatusOK, []byte("....")
}
