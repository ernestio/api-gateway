package helpers

// User : interface for users
type User interface {
	GetAdmin() bool
	GetGroupID() int
}

var (
	// AuthNonAdmin : Response body for non authorized requests on admin resources
	AuthNonAdmin = []byte("You don't have permissions to perform this action, please login with an admin account")
	// AuthNonOwner : Response body for non authorized requests on owned resources
	AuthNonOwner = []byte("You don't have permissions to perform this action, please login as a resource owner")
	// AuthNonReadable : Response body for non authorized requests on admin resources
	AuthNonReadable = []byte("You don't have permissions to perform this action, please contact the resource owner")
	// AuthNonGroup : Response body for non autherized requests due to a non group users
	AuthNonGroup = []byte("Current user does not belong to any group.\nPlease assign the user to a group before performing this action")
)

// IsAuthorized : Validates if the given user has access to the given resource
func IsAuthorized(au User, resource string) (int, []byte) {
	resourceID := ""

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

	if au.GetAdmin() == true {
		return 200, []byte("")
	}

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

	groupResources := map[string]int{
		"services/create": 401,
	}
	if st, ok := groupResources[resource]; ok {
		if au.GetGroupID() == 0 {
			return st, AuthNonGroup
		}
	}

	if resourceID != "" {
		ownedResources := map[string]int{}
		if st, ok := ownedResources[resource]; ok {
			if !IsOwner(au, resource, resourceID) {
				return st, AuthNonOwner
			}
		}

		readableResources := map[string]int{}
		if st, ok := readableResources[resource]; ok {
			if !IsReader(au, resource, resourceID) {
				return st, AuthNonReadable
			}
		}
	}

	return 200, []byte("")
}

// IsOwner : Checks if the given user is owner of a specific resource
func IsOwner(au User, resource, resourceID string) bool {
	return true
}

// IsReader : Checks if the given user has read permissions on a specific resource
func IsReader(au User, resource, resourceID string) bool {
	return true
}
