package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngin8-beta/tfbackend/internal/storage"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	storage, err := GetStorage()
	if err != nil {
		log.Fatalf("failed to get storage: %v", err)
	}

	r := gin.Default()
	r.GET("/:project/state", GetStateHundler(storage))
	r.POST("/:project/state", PostStateHundler(storage))
	r.Handle("LOCK", "/:project/state", LockStateHundler)
	r.Handle("UNLOCK", "/:project/state", UnlockStateHundler)
	r.NoRoute(NoRouteStateHundler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	return &Server{httpServer: srv}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func GetStateHundler(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := storage.GetState(c.Param("project"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
	}
}

func PostStateHundler(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var state map[string]interface{}
		if err := c.ShouldBindJSON(&state); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		storage.PostState(c.Param("project"), state)
	}
}

func LockStateHundler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "lock state",
	})
}

func UnlockStateHundler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "unlock state",
	})
}

func NoRouteStateHundler(c *gin.Context) {
	c.JSON(404, gin.H{
		"message": "not found",
	})
}