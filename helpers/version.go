package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/blang/semver"
)

const RequiredCliVersion string = "2.2.0"

// ValidCliVersion checks to see if the client version meets minimum
// requirements
func ValidCliVersion(r *http.Request) error {
	for _, v := range r.Header["User-Agent"] {
		if strings.Contains(v, "Ernest/") {
			parts := strings.Split(v, "/")
			ernestVersion := parts[1]

			rv, err := semver.Make(RequiredCliVersion)
			if err != nil {
				return err
			}
			ev, err := semver.Make(ernestVersion)
			if err != nil {
				return err
			}
			if ev.LT(rv) {
				err := fmt.Sprintf("Ernest CLI %s is not supported by this server.\nPlease upgrade http://docs.ernest.io/downloads/", ernestVersion)
				return errors.New(err)
			}
		}
	}
	return nil
}
