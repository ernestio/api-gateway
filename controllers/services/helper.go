package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

	h "github.com/ernestio/api-gateway/helpers"
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

func generateStreamID(salt string) string {
	compose := []byte(salt)
	hasher := md5.New()
	if _, err := hasher.Write(compose); err != nil {
		h.L.Warning(err.Error())
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// Generates a service id composed by a random uuid, and
// a valid generated stream id
func generateServiceID(salt string) string {
	sufix := generateStreamID(salt)
	prefix, _ := uuid.NewV4()

	return prefix.String() + "-" + string(sufix[:])
}
