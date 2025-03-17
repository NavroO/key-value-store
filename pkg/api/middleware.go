package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) logMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start)

	s.logRequest(c.Request, duration)
}

func (s *Server) logRequest(req *http.Request, duration time.Duration) {
	s.logger.Info("Request",
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("query", req.URL.RawQuery),
		zap.String("remote", req.RemoteAddr),
		zap.Duration("duration", duration),
	)
}
