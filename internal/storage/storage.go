package storage

type Storage interface {
	// GetState retrieves the state file corresponding to the specified project.
	GetState(project string) (map[string]interface{}, error)
	// PostState saves the state as a JSON file for the specified project.
	PostState(project string, state map[string]interface{}) error
}

func MinimalState() map[string]interface{} {
	return map[string]interface{}{
		"version": 1,
	}
}