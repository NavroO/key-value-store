package api

import (
	"context"
	"net/http"
	"os"

	"github.com/NavroO/go-key-value-store/pkg/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	store  *storage.InMemoryStorage
	router *gin.Engine
	logger *zap.Logger
	server *http.Server
}

func NewServer(store *storage.InMemoryStorage) *Server {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		logFilePath = "logs.txt"
	}

	addr := os.Getenv("APP_PORT")
	if addr == "" {
		addr = ":8080"
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter := zapcore.AddSync(file)
	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zap.InfoLevel),
		zapcore.NewCore(encoder, consoleWriter, zap.InfoLevel),
	)

	logger := zap.New(core)

	server := &Server{
		store:  store,
		router: gin.Default(),
		logger: logger,
		server: &http.Server{Addr: addr, Handler: router},
	}

	server.router.Use(server.logMiddleware)
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/store", s.setKey)
	s.router.GET("/store/:key", s.getKey)
	s.router.DELETE("/store/:key", s.deleteKey)
}

func (s *Server) Run(port string) {
	s.server.Addr = port
	s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
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
