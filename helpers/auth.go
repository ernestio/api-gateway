package helpers

// User : interface for users
type User interface {
	GetAdmin() bool
}

var (
	// AuthNonAdmin : Response body for non authorized requests on admin resources
	AuthNonAdmin = []byte("You don't have permissions to perform this action, please login with an admin account")
)

// IsAuthorized : Validates if the given user has access to the given resource
func IsAuthorized(au User, resource string) (int, []byte) {
	if au.GetAdmin() == false {
		adminResources := map[string]int{
			"groups/add_datacenter":     403,
			"groups/add_user":           403,
			"groups/create":             403,
			"groups/delete":             403,
			"groups/rm_datacenter":      403,
			"groups/rm_user":            403,
			"groups/update":             403,
			"loggers/create":            403,
			"loggers/delete":            403,
			"loggers/list":              403,
			"notifications/add_service": 403,
			"notifications/create":      403,
			"notifications/delete":      403,
			"notifications/list":        403,
			"notifications/rm_service":  403,
			"notifications/update":      403,
			"usages/report":             403,
			"users/create":              403,
			"users/delete":              403,
		}
		if st, ok := adminResources[resource]; ok {
			return st, AuthNonAdmin
		}
	}

	licensedResources := map[string]int{
		"notifications/add_service": 405,
		"notifications/create":      405,
		"notifications/delete":      405,
		"notifications/list":        405,
		"notifications/rm_service":  405,
		"notifications/update":      405,
		"services/sync":             405,
		"services/update":           405,
	}
	if st, ok := licensedResources[resource]; ok {
		if err := Licensed(); err != nil {
			return st, []byte(err.Error())
		}
	}

	return 200, []byte("")
}
