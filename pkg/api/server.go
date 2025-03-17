package api

import (
	"net/http"

	"github.com/NavroO/go-key-value-store/pkg/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *storage.InMemoryStorage
	router *gin.Engine
}

func NewServer(store *storage.InMemoryStorage) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/store", s.setKey)
	s.router.GET("/store/:key", s.getKey)
	s.router.DELETE("/store/:key", s.deleteKey)
}

func (s *Server) Run(port string) {
	s.router.Run(":" + port)
}

func (s *Server) setKey(c *gin.Context) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Ttl   int    `json:"ttl"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	s.store.Set(req.Key, req.Value, req.Ttl)
	c.JSON(http.StatusOK, gin.H{"message": "Value set successfully"})
}

func (s *Server) getKey(c *gin.Context) {
	key := c.Param("key")
	value, exists := s.store.Get(key)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func (s *Server) deleteKey(c *gin.Context) {
	key := c.Param("key")
	deleted := s.store.Delete(key)
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key deleted"})
}
