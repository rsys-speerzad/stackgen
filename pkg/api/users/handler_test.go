package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (m *mockStore) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *mockStore) Get(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *mockStore) Update(id string, user *models.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}
func (m *mockStore) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockStore) AddAvailability(slots []models.UserAvailability) error {
	args := m.Called(slots)
	return args.Error(0)
}
func (m *mockStore) UpdateAvailability(id string, slot models.UserAvailability) error {
	args := m.Called(id, slot)
	return args.Error(0)
}
func (m *mockStore) DeleteAvailability(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *mockStore) GetAvailability(userID string) ([]models.UserAvailability, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserAvailability), args.Error(1)
}

func newHandlerWithMockStore(store *mockStore) *Handler {
	h := &Handler{store: store}
	return h
}

func TestCreate_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	store.On("Create", user).Return(nil)
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusCreated, w.Code)
	store.AssertExpectations(t)
}

func TestCreate_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	store.On("Create", user).Return(errors.New("fail"))
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	store.On("Get", "1").Return(user, nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestGet_NotFound(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Get", "1").Return(&models.User{}, gorm.ErrRecordNotFound)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusNotFound, w.Code)
	store.AssertExpectations(t)
}

func TestGet_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Get", "1").Return(&models.User{}, errors.New("fail"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestGet_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/", nil)
	params := httprouter.Params{}
	h.Get(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	store.On("Update", "1", user).Return(nil)
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Update(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestUpdate_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Update(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	store.On("Update", "1", user).Return(errors.New("fail"))
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Update(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestUpdate_BadRequest_NoID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	user := &models.User{ID: uuid.New()}
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPut, "/users/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	params := httprouter.Params{}
	h.Update(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Delete", "1").Return(nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusNoContent, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Delete", "1").Return(gorm.ErrRecordNotFound)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusNotFound, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	store.On("Delete", "1").Return(errors.New("fail"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	params := httprouter.Params{{Key: "id", Value: "1"}}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestDelete_BadRequest(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/", nil)
	params := httprouter.Params{}
	h.Delete(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAvailability_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	now := time.Now()
	slots := []models.Slot{
		{StartTime: now, EndTime: now.Add(30 * time.Minute)},
		{StartTime: now.Add(60 * time.Minute), EndTime: now.Add(90 * time.Minute)},
	}
	reqBody, _ := json.Marshal(map[string]interface{}{
		"slots": slots,
	})
	// var userAvailabilities []models.UserAvailability
	// for _, slot := range slots {
	// 	userAvailabilities = append(userAvailabilities, models.UserAvailability{
	// 		UserID: userID,
	// 		Slot:   slot,
	// 	})
	// }
	store.On("AddAvailability", mock.AnythingOfType("[]models.UserAvailability")).Return(nil).Run(func(args mock.Arguments) {
		// Optionally check the slots
	})
	r := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/availability", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}}
	h.AddAvailability(w, r, params)
	assert.Equal(t, http.StatusNoContent, w.Code)
	store.AssertExpectations(t)
}

func TestAddAvailability_BadRequest_NoUserID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPost, "/users//availability", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	params := httprouter.Params{}
	h.AddAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAvailability_BadRequest_InvalidUserID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPost, "/users/invalid-uuid/availability", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: "invalid-uuid"}}
	h.AddAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAvailability_BadRequest_InvalidBody(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/availability", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}}
	h.AddAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAvailability_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	now := time.Now()
	slots := []models.Slot{
		{StartTime: now, EndTime: now.Add(30 * time.Minute)},
		{StartTime: now.Add(60 * time.Minute), EndTime: now.Add(90 * time.Minute)},
	}
	reqBody, _ := json.Marshal(map[string]interface{}{"slots": slots})
	store.On("AddAvailability", mock.AnythingOfType("[]models.UserAvailability")).Return(errors.New("fail"))
	r := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/availability", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}}
	h.AddAvailability(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestUpdateAvailability_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	aid := "avail-1"
	now := time.Now()
	slot := models.UserAvailability{Slot: models.Slot{StartTime: now, EndTime: now.Add(30 * time.Minute)}}
	body, _ := json.Marshal(slot)
	store.On("UpdateAvailability", aid, mock.AnythingOfType("models.UserAvailability")).Return(nil)
	r := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/availability/"+aid, bytes.NewReader(body))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}, {Key: "aid", Value: aid}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusNoContent, w.Code)
	store.AssertExpectations(t)
}

func TestUpdateAvailability_BadRequest_NoUserID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPut, "/users//availability/aid", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "aid", Value: "aid"}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAvailability_BadRequest_InvalidUserID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodPut, "/users/invalid-uuid/availability/aid", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: "invalid-uuid"}, {Key: "aid", Value: "aid"}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAvailability_BadRequest_NoAid(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	r := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/availability/", bytes.NewReader([]byte("{}")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAvailability_BadRequest_InvalidBody(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	aid := "aid"
	r := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/availability/"+aid, bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}, {Key: "aid", Value: aid}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAvailability_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := uuid.New()
	aid := "aid"
	now := time.Now()
	slot := models.UserAvailability{Slot: models.Slot{StartTime: now, EndTime: now.Add(30 * time.Minute)}}
	body, _ := json.Marshal(slot)
	store.On("UpdateAvailability", aid, mock.AnythingOfType("models.UserAvailability")).Return(errors.New("fail"))
	r := httptest.NewRequest(http.MethodPut, "/users/"+userID.String()+"/availability/"+aid, bytes.NewReader(body))
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID.String()}, {Key: "aid", Value: aid}}
	h.UpdateAvailability(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestDeleteAvailability_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	aid := "aid"
	store.On("DeleteAvailability", aid).Return(nil)
	r := httptest.NewRequest(http.MethodDelete, "/users/1/availability/"+aid, nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "aid", Value: aid}}
	h.DeleteAvailability(w, r, params)
	assert.Equal(t, http.StatusNoContent, w.Code)
	store.AssertExpectations(t)
}

func TestDeleteAvailability_NotFound(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	aid := "aid"
	store.On("DeleteAvailability", aid).Return(gorm.ErrRecordNotFound)
	r := httptest.NewRequest(http.MethodDelete, "/users/1/availability/"+aid, nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "aid", Value: aid}}
	h.DeleteAvailability(w, r, params)
	assert.Equal(t, http.StatusNotFound, w.Code)
	store.AssertExpectations(t)
}

func TestDeleteAvailability_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	aid := "aid"
	store.On("DeleteAvailability", aid).Return(errors.New("fail"))
	r := httptest.NewRequest(http.MethodDelete, "/users/1/availability/"+aid, nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "aid", Value: aid}}
	h.DeleteAvailability(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}

func TestDeleteAvailability_BadRequest_NoAid(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodDelete, "/users/1/availability/", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{}
	h.DeleteAvailability(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAvailabilities_Success(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := "user-1"
	now := time.Now()
	availabilities := []models.UserAvailability{
		{Slot: models.Slot{StartTime: now, EndTime: now.Add(30 * time.Minute)}},
		{Slot: models.Slot{StartTime: now.Add(60 * time.Minute), EndTime: now.Add(90 * time.Minute)}},
	}
	store.On("GetAvailability", userID).Return(availabilities, nil)
	r := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/availability", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID}}
	h.GetAvailabilities(w, r, params)
	assert.Equal(t, http.StatusOK, w.Code)
	store.AssertExpectations(t)
}

func TestGetAvailabilities_BadRequest_NoUserID(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	r := httptest.NewRequest(http.MethodGet, "/users//availability", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{}
	h.GetAvailabilities(w, r, params)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAvailabilities_Error(t *testing.T) {
	store := new(mockStore)
	h := newHandlerWithMockStore(store)
	userID := "user-1"
	store.On("GetAvailability", userID).Return([]models.UserAvailability{}, errors.New("fail"))
	r := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/availability", nil)
	w := httptest.NewRecorder()
	params := httprouter.Params{{Key: "id", Value: userID}}
	h.GetAvailabilities(w, r, params)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	store.AssertExpectations(t)
}
