package storage

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return fmt.Errorf("failed to create directory: %w", err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state to JSON: %w", err)
	}

	// Backup the existing file if it exists.
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Rename(filePath, filePath+".backup"); err != nil {
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	// Writes to a temporary file and performs atomic updates.
	tmpFilePath := filePath + ".tmp"
	if err := os.WriteFile(tmpFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := os.Rename(tmpFilePath, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	return nil
}

// getFilePath constructs the file path for saving based on the project name.
// Example: <basePath>/<project>/terraform.tfstate
func (ls *LocalStorage) getFilePath(project string) string {
	return filepath.Join(ls.basePath, project, "terraform.tfstate")
}
