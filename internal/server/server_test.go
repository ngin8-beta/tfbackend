package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngin8-beta/tfbackend/internal/server"
)

func NewGinServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetStateHundler(t *testing.T) {
	// Arrange
	r := NewGinServer()
	r.GET("/myProject/state", server.GetStateHundler)
	req, _ := http.NewRequest("GET", "/myProject/state", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

func TestPostStateHundler(t *testing.T) {
	// Arrange
	r := NewGinServer()
	r.POST("/myProject/state", server.PostStateHundler)
	req, _ := http.NewRequest("POST", "/myProject/state", nil)
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}

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