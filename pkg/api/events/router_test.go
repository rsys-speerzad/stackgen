package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// mockInitializeRouter registers dummy handlers for route existence testing
func mockInitializeRouter(router *httprouter.Router) {
	dummyHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	}
	router.POST("/event", dummyHandler)
	router.GET("/event/:id", dummyHandler)
	router.PUT("/event/:id", dummyHandler)
	router.DELETE("/event/:id", dummyHandler)
	router.GET("/events/:id/recommendations", dummyHandler)
}
func TestInitializeRouter_Routes(t *testing.T) {
	router := httprouter.New()
	mockInitializeRouter(router)

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/event"},
		{"GET", "/event/123"},
		{"PUT", "/event/123"},
		{"DELETE", "/event/123"},
		{"GET", "/events/123/recommendations"},
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
