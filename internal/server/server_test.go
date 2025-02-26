package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngin8-beta/tfbackend/internal/server"
	"github.com/ngin8-beta/tfbackend/internal/storage"
)

func NewGinServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetStateHundler(t *testing.T) {
	baseDir := t.TempDir()
	storage, _ := storage.NewLocalStorage(baseDir)
	prjDir := baseDir + "/myProject"
	os.Mkdir(prjDir, 0777)
	os.WriteFile(prjDir+"/state.json", []byte(`{"key":"value"}`), 0666)

	r := gin.Default()
	r.GET("/:project/state", server.GetStateHundler(storage))
	req, _ := http.NewRequest("GET", "/myProject/state", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

// func TestPostStateHundler(t *testing.T) {
// 	// Arrange
// 	r := NewGinServer()
// 	storage, _ := server.GetStorage()
// 	r.POST("/myProject/state", server.PostStateHundler(storage))
// 	req, _ := http.NewRequest("POST", "/myProject/state", nil)
// 	w := httptest.NewRecorder()
// 
// 	// Act
// 	r.ServeHTTP(w, req)
// 	fmt.Println("w.Body: ", w.Body)
// 	fmt.Println("w.Body.String(): ", w.Body.String())
// 	fmt.Println(req.Body)
// 
// 	// Assert
// 	if w.Code != http.StatusOK {
// 		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
// 	}
// }

func TestLockStateHundler(t *testing.T) {
	// Arrange
	r := NewGinServer()
	r.Handle("LOCK", "/:project/state", server.LockStateHundler)
	req, _ := http.NewRequest("LOCK", "/myProject/state", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

func TestUnlockStateHundler(t *testing.T) {
	// Arrange
	r := NewGinServer()
	r.Handle("UNLOCK", "/:project/state", server.UnlockStateHundler)
	req, _ := http.NewRequest("UNLOCK", "/myProject/state", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

func TestNoRouteStateHundler(t *testing.T) {
	// Arrange
	r := NewGinServer()
	r.NoRoute(server.NoRouteStateHundler)
	req, _ := http.NewRequest("GET", "/myProject/state", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
	}
}