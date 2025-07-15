package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rsys-speerzad/stackgen/pkg/models"
	"github.com/rsys-speerzad/stackgen/pkg/router"
	"github.com/rsys-speerzad/stackgen/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestUserEventRecommendationFlow(t *testing.T) {
	// setup database connection
	// (Assuming store.GetDB() is set up to use the test database)
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASS", "admin")
	os.Setenv("DB_NAME", "stackgen")
	store.InitDB()

	// create a test server
	server := httptest.NewServer(router.NewServer().Handler)
	defer server.Close()

	// 1. Create a user
	user := models.User{Name: "Alice", Email: "alice@example.com"}
	userBody, _ := json.Marshal(user)
	resp, err := http.Post(server.URL+"/user", "application/json", bytes.NewReader(userBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var createdUser models.User
	json.NewDecoder(resp.Body).Decode(&createdUser)
	resp.Body.Close()

	defer func() {
		// Cleanup: delete the created user after the test
		req := httptest.NewRequest(http.MethodDelete, server.URL+"/user/"+createdUser.ID.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		http.DefaultClient.Do(req)
	}()

	// 2. Add availability for the user
	now := time.Now()
	avail := struct {
		Slots []models.Slot `json:"slots"`
	}{
		Slots: []models.Slot{
			{StartTime: now, EndTime: now.Add(30 * time.Minute)},
			{StartTime: now.Add(60 * time.Minute), EndTime: now.Add(90 * time.Minute)},
			{StartTime: now.Add(120 * time.Minute), EndTime: now.Add(150 * time.Minute)},
			{StartTime: now.Add(180 * time.Minute), EndTime: now.Add(210 * time.Minute)},
		},
	}
	availBody, _ := json.Marshal(avail)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/user/"+createdUser.ID.String()+"/availability", bytes.NewReader(availBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 3. Create an event
	event := models.Event{Title: "Team Meeting", EstimatedDuration: 60, EventSlots: []models.EventSlot{
		{StartTime: now, EndTime: now.Add(60 * time.Minute)},
		{StartTime: now.Add(120 * time.Minute), EndTime: now.Add(180 * time.Minute)},
	}}
	eventBody, _ := json.Marshal(event)
	resp, err = http.Post(server.URL+"/event", "application/json", bytes.NewReader(eventBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var createdEvent models.Event
	json.NewDecoder(resp.Body).Decode(&createdEvent)
	resp.Body.Close()

	// 4. Get recommendations
	resp, err = http.Get(server.URL + "/events/" + createdEvent.ID.String() + "/recommendations")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Optionally decode and check recommendations here
	resp.Body.Close()
}
