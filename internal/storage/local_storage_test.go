package storage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ngin8-beta/tfbackend/internal/storage"
)

func TestNewLocalStorageWithEmptyBasePath(t *testing.T) {
	_, err := storage.NewLocalStorage("")
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
}

func TestNewLocalStorageWithNonExistentBasePath(t *testing.T) {
	basePath := "/tmp/test"
	_, err := storage.NewLocalStorage(basePath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	os.Remove(basePath)
}

func TestNewLocalStorageWithExistentBasePath(t *testing.T) {
	basePath := "/tmp/test"
	os.Mkdir(basePath, 0755)
	_, err := storage.NewLocalStorage(basePath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	os.Remove(basePath)
}

func TestNewLocalStorageWithNonWritableBasePath(t *testing.T) {
	basePath := "/tmp/test"
	os.Mkdir(basePath, 0444)
	_, err := storage.NewLocalStorage(basePath + "/test")
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	os.Remove(basePath)
}

func TestGetState(t *testing.T) {
	basePath := t.TempDir()
	ls, err := storage.NewLocalStorage(basePath)
	if err != nil {
		t.Fatalf("failed to initialize LocalStorage: %v", err)
	}

	project := "testproject"
	state := map[string]interface{}{"foo": "bar"}
	filePath := filepath.Join(basePath, project, ".json")
	dir := filepath.Dir(filePath)
	os.Mkdir(dir, 0755)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}
	os.WriteFile(filePath, data, 0644)

	readState, err := ls.GetState(project)
	if err != nil {
		t.Errorf("GetState returned unexpected error: %v", err)
	}

	if readState == nil {
		t.Fatalf("GetState returned nil state")
	}

	if readState["foo"] != "bar" {
		t.Errorf("expected 'foo': 'bar', got: %v", readState)
	}
}

func TestPostState(t *testing.T) {
	basePath := t.TempDir()
	ls, err := storage.NewLocalStorage(basePath)
	if err != nil {
		t.Fatalf("failed to initialize LocalStorage: %v", err)
	}

	project := "testproject"
	state := map[string]interface{}{"foo": "bar"}
	err = ls.PostState(project, state)
	if err != nil {
		t.Errorf("PostState returned unexpected error: %v", err)
	}

	filePath := filepath.Join(basePath, project, ".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist, but it does not", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read state file: %v", err)
	}
	if !contains(content, `"foo": "bar"`) {
		t.Errorf("expected 'foo': 'bar' in file content, got: %s", string(content))
	}
}

// helper to check content quickly
func contains(data []byte, str string) bool {
	return string(data) == str || len(data) > len(str) && string(data[:len(str)]) == str || string(data[len(data)-len(str):]) == str || string(data) != ""
}


