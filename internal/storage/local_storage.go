package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// LocalStorage is a storage implementation that uses the local file system.
type LocalStorage struct {
	// The root directory for saving state files.
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if basePath == "" {
		return nil, errors.New("basePath cannot be empty")
	}
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: basePath}, nil
}

// getFilePath constructs the file path for saving based on the project name.
// Example: <basePath>/<project>/state.json
func (ls *LocalStorage) getFilePath(project string) string {
	return filepath.Join(ls.basePath, project, "state.json")
}

func (ls *LocalStorage) GetName() string {
	return "LocalFileSystem"
}

func (ls *LocalStorage) GetState(project string) (map[string]interface{}, error) {
	filePath := ls.getFilePath(project)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return MinimalState(), nil
		}
		return nil, err
	}
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func (ls *LocalStorage) PostState(project string, state map[string]interface{}) error {
	filePath := ls.getFilePath(project)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func (ls *LocalStorage) DeleteState(project string) error {
	filePath := ls.getFilePath(project)
	return os.Remove(filePath)
}

