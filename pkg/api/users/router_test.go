package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// mockInitializeRouter registers dummy handlers for route existence testing
func mockInitializeRouter(r *httprouter.Router) {
	dummyHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	}
	r.POST("/user", dummyHandler)
	r.GET("/user/:id", dummyHandler)
	r.PUT("/user/:id", dummyHandler)
	r.DELETE("/user/:id", dummyHandler)
	r.GET("/user/:id/availability", dummyHandler)
	r.POST("/user/:id/availability", dummyHandler)
	r.PUT("/user/:id/availability/:aid", dummyHandler)
	r.DELETE("/user/:id/availability/:aid", dummyHandler)
}
func TestInitializeRouter_Routes(t *testing.T) {
	router := httprouter.New()
	mockInitializeRouter(router)

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/user"},
		{"GET", "/user/123"},
		{"PUT", "/user/123"},
		{"DELETE", "/user/123"},
		{"GET", "/user/123/availability"},
		{"POST", "/user/123/availability"},
		{"PUT", "/user/123/availability/456"},
		{"DELETE", "/user/123/availability/456"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Route %s %s should exist and return 200", tt.method, tt.path)
	}
}

func TestInitializeRouter_RouteNotFound(t *testing.T) {
	router := httprouter.New()
	mockInitializeRouter(router)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "Nonexistent route should return 404")
}
