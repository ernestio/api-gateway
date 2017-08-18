package helpers

// User : interface for users
type User interface {
	GetAdmin() bool
	IsOwner(resourceType, resourceID string) bool
	IsReader(resourceType, resourceID string) bool
}

var (
	// AuthNonAdmin : Response body for non authorized requests on admin resources
	AuthNonAdmin = []byte("You don't have permissions to perform this action, please login with an admin account")
	// AuthNonOwner : Response body for non authorized requests on owned resources
	AuthNonOwner = []byte("You don't have permissions to perform this action, please login as a resource owner")
	// AuthNonReadable : Response body for non authorized requests on admin resources
	AuthNonReadable = []byte("You don't have permissions to perform this action, please contact the resource owner")
	// GetProject : ...
	GetProject = "get_project"
	// DeleteProject : ...
	DeleteProject = "delete_project"
	// UpdateProject : ...
	UpdateProject = "update_project"
	// DeleteEnv : ...
	DeleteEnv = "delete_env"
	// DeleteEnvForce : ..
	DeleteEnvForce = "delete_env_force"
	// UpdateEnv : ...
	UpdateEnv = "update_env"
	// GetEnv : ...
	GetEnv = "get_environment"
	// SyncEnv : ...
	SyncEnv = "sync_env"
	// ListBuilds : ...
	ListBuilds = "list_builds"
	// DeleteBuild : ...
	DeleteBuild = "delete_build"
	// GetBuild : ...
	GetBuild = "get_build"
	// ResetBuild : ...
	ResetBuild = "reset_build"
)

// IsAuthorized : Validates if the given user has access to the given resource
func IsAuthorized(au User, resource string) (int, []byte) {
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

	return 200, []byte("")
}

// IsAuthorizedToResource : check  if the user is authorized to access a specific resource
func IsAuthorizedToResource(au User, endpoint, resource, resourceID string) (int, []byte) {
	if resourceID != "" {
		ownedResources := map[string]int{
			DeleteBuild:    403,
			DeleteEnv:      403,
			DeleteEnvForce: 403,
			UpdateEnv:      403,
			DeleteProject:  403,
			UpdateProject:  403,
			ResetBuild:     403,
			SyncEnv:        403,
		}
		if st, ok := ownedResources[endpoint]; ok {
			if !au.IsOwner(resource, resourceID) {
				return st, AuthNonOwner
			}
		}

		readableResources := map[string]int{
			GetProject: 403,
			GetEnv:     403,
			ListBuilds: 403,
			GetBuild:   403,
		}
		if st, ok := readableResources[endpoint]; ok {
			if !au.IsReader(resource, resourceID) {
				return st, AuthNonReadable
			}
		}
	}

	return 200, []byte("")
}
