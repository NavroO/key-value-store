package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NavroO/go-key-value-store/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestServer() (*Server, *storage.InMemoryStorage) {
	store := storage.NewInMemoryStorage()
	server := NewServer(store)
	return server, store
}

func TestMainAPI(t *testing.T) {
	t.Run("Set Key", func(t *testing.T) {
		server, _ := setupTestServer()

		payload := []byte(`{"key": "username", "value": "joe"}`)
		req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		server.router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Get Key", func(t *testing.T) {
		server, store := setupTestServer()

		store.Set("username", "joe", 0)

		req, _ := http.NewRequest("GET", "/store/username", nil)
		resp := httptest.NewRecorder()
		server.router.ServeHTTP(resp, req)

		expectedResponse := `{"key":"username","value":"joe"}`

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, expectedResponse, resp.Body.String())
	})

	t.Run("Delete Key", func(t *testing.T) {
		server, store := setupTestServer()

		store.Set("username", "joe", 0)

		req, _ := http.NewRequest("DELETE", "/store/username", nil)
		resp := httptest.NewRecorder()
		server.router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		_, exists := store.Get("username")
		assert.False(t, exists, "Klucz powinien zostać usunięty")
	})

	t.Run("Get Non-Existent Key", func(t *testing.T) {
		server, _ := setupTestServer()

		req, _ := http.NewRequest("GET", "/store/nonexistent", nil)
		resp := httptest.NewRecorder()
		server.router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error":"Key not found"}`, resp.Body.String())
	})

	t.Run("Set Key with TTL", func(t *testing.T) {
		server, store := setupTestServer()

		payload := []byte(`{"key": "session", "value": "abcdef", "ttl": 2}`)
		req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		server.router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		_, exists := store.Get("session")
		assert.True(t, exists, "Klucz powinien istnieć zaraz po dodaniu")

		time.Sleep(3 * time.Second)
		_, exists = store.Get("session")
		assert.False(t, exists, "Klucz powinien wygasnąć po TTL")
	})
}
