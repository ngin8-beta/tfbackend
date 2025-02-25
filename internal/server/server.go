package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string) *Server {
	r := gin.Default()
	r.GET("/:project/state", GetStateHundler)
	r.POST("/:project/state", PostStateHundler)
	r.Handle("LOCK", "/:project/state", LockStateHundler)
	r.Handle("UNLOCK", "/:project/state", UnlockStateHundler)
	r.NoRoute(NoRouteStateHundler)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &Server{httpServer: srv}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func GetStateHundler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get state",
	})
}

func PostStateHundler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "post state",
	})
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
