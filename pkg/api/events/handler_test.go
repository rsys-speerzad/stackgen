package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/rsys-speerzad/stackgen/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Create(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
func (m *mockStore) Get(id string) (*models.Event, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Event), args.Error(1)
}
func (m *mockStore) Update(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
func (m *mockStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockStore) GetRecommendations(eventID string) (*models.RecommendedSlot, error) {
	args := m.Called(eventID)
	return args.Get(0).(*models.RecommendedSlot), args.Error(1)
}

func newHandlerWithMockStore(store *mockStore) *Handler {
	h := &Handler{store: store}
	return h
}

func TestCreate_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	event := &models.Event{ID: uuid.New(), Title: "Test"}
	store.On("Create", event).Return(nil)
	body, _ := json.Marshal(event)
	r := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusCreated, w.Code)
	store.AssertExpectations(t)
}

func TestCreate_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	event := &models.Event{ID: uuid.New(), Title: "Test"}
	store.On("Create", event).Return(errors.New("fail"))
	body, _ := json.Marshal(event)
	r := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	id := uuid.New()
	event := &models.Event{ID: id, Title: "Test"}
	store.On("Get", id.String()).Return(event, nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events/%s", id.String()), nil)
	params := httprouter.Params{{Key: "id", Value: id.String()}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestGet_NotFound(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Get", "1").Return(&models.Event{}, gorm.ErrRecordNotFound)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/events/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusNotFound, w.Code)
	store.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Get", "1").Return(&models.Event{}, errors.New("fail"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/events/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestGet_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/events/", nil)
	params := httprouter.Params{}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	id := uuid.New()
	event := &models.Event{ID: id, Title: "Test"}
	store.On("Update", event).Return(nil)
	body, _ := json.Marshal(event)
	r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/events/%s", id.String()), bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Update(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestUpdate_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPut, "/events/1", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	h.Update(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	id := uuid.New()
	event := &models.Event{ID: id, Title: "Test"}
	store.On("Update", event).Return(errors.New("fail"))
	body, _ := json.Marshal(event)
	r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/events/%s", id.String()), bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Update(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Delete", "1").Return(nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	id := uuid.New()
	store.On("Delete", id.String()).Return(gorm.ErrRecordNotFound)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/events/%s", id.String()), nil)
	params := httprouter.Params{{Key: "id", Value: id.String()}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusNotFound, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	id := uuid.New()
	store.On("Delete", id.String()).Return(errors.New("fail"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/events/%s", id.String()), nil)
	params := httprouter.Params{{Key: "id", Value: id.String()}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/events/", nil)
	params := httprouter.Params{}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRecommendations_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	recs := &models.RecommendedSlot{}
	id := uuid.New()
	store.On("GetRecommendations", id.String()).Return(recs, nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events/%s/recommendations", id.String()), nil)
	params := httprouter.Params{{Key: "id", Value: id.String()}}
	h.GetRecommendations(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestGetRecommendations_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/events//recommendations", nil)
	params := httprouter.Params{}
	h.GetRecommendations(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetRecommendations_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("GetRecommendations", "1").Return(&models.RecommendedSlot{}, errors.New("fail"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/events/1/recommendations", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.GetRecommendations(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}
