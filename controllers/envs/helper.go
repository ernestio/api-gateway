package envs

import (
	"encoding/json"
	"errors"

	"github.com/ernestio/api-gateway/models"
	"github.com/nu7hatch/gouuid"
)

func getEnvRaw(au models.User, name string) ([]byte, error) {
	filters := make(map[string]interface{}, 0)
	filters["name"] = name

	ss, err := au.EnvsBy(filters)
	if err != nil {
		return nil, err
	}

	if len(ss) == 0 {
		return nil, errors.New("Not found")
	}

	body, err := json.Marshal(ss[0])
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

// Generates an environment id composed by a random uuid, and
// a valid generated stream id
func generateEnvID(salt string) string {
	id, _ := uuid.NewV4()

	return id.String()
}
