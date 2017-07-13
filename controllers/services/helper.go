package services

import (
	"encoding/json"
	"errors"

	"github.com/ernestio/api-gateway/models"
	"github.com/nu7hatch/gouuid"
)

func getServiceRaw(name string, group int) ([]byte, error) {
	var ss models.Service

	s, err := ss.GetByNameAndGroupID(name, group)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(s)
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

// Generates a service id composed by a random uuid, and
// a valid generated stream id
func generateServiceID(salt string) string {
	id, _ := uuid.NewV4()

	return id.String()
}
