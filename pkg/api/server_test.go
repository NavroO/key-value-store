package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NavroO/go-key-value-store/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSetKey(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewServer(store)

	payload := []byte(`{"key": "username", "value": "joe"}`)
	req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetKey(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewServer(store)

	store.Set("username", "joe")

	req, _ := http.NewRequest("GET", "/store/username", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	expectedResponse := `{"key":"username","value":"joe"}`

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, expectedResponse, resp.Body.String())
}

func TestDeleteKey(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewServer(store)

	store.Set("username", "joe")

	req, _ := http.NewRequest("DELETE", "/store/username", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	_, exists := store.Get("username")
	assert.False(t, exists, "Klucz powinien zostać usunięty")
}

func TestGetNonExistentKey(t *testing.T) {
	store := storage.NewInMemoryStorage()
	server := NewServer(store)

	req, _ := http.NewRequest("GET", "/store/nonexistent", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.JSONEq(t, `{"error":"Key not found"}`, resp.Body.String())
}
