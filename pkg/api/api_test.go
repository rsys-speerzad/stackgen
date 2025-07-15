package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rsys-speerzad/stackgen/pkg/models"
)

func TestSuccessJson_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	data := map[string]string{"foo": "bar"}

	SuccessJson(rr, req, data)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `"foo":"bar"`) {
		t.Errorf("expected body to contain foo:bar, got %s", string(body))
	}
}

func TestSuccessJson_MarshalError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ch := make(chan int) // not serializable

	SuccessJson(rr, req, ch)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "serialising response failed") {
		t.Errorf("expected marshal error, got %s", string(body))
	}
}

func TestSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo", nil)
	msg := []byte(`{"ok":true}`)

	Success(rr, req, msg)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ao := resp.Header.Get("Access-Control-Allow-Origin"); ao != "*" {
		t.Errorf("expected Access-Control-Allow-Origin *, got %s", ao)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != string(msg) {
		t.Errorf("expected body %s, got %s", msg, body)
	}
}

func TestError_DefaultStatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/err", nil)
	err := models.ErrNotFound

	Error(rr, req, err, 0)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), err.Error()) {
		t.Errorf("expected error message in body, got %s", string(body))
	}
}

func TestError_CustomStatusCode(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/err", nil)
	err := errors.New("custom error")

	Error(rr, req, err, http.StatusTeapot)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("expected status 418, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "custom error") {
		t.Errorf("expected error message in body, got %s", string(body))
	}
}

func TestError_NilError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/err", nil)

	Error(rr, req, nil, 400)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "nil err") {
		t.Errorf("expected 'nil err' in body, got %s", string(body))
	}
}

func TestToHTTPStatusCode(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{models.ErrMissingArgument, http.StatusBadRequest},
		{models.ErrInvalidMessageType, http.StatusBadRequest},
		{models.ErrNotFound, http.StatusNotFound},
		{errors.New("other"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		got := toHTTPStatusCode(tt.err)
		if got != tt.expected {
			t.Errorf("toHTTPStatusCode(%v) = %d, want %d", tt.err, got, tt.expected)
		}
	}
}

func TestResponseWriter_StatusOK_WithData(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}

	ResponseWriter(rr, data, 0)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `"hello":"world"`) {
		t.Errorf("expected body to contain hello:world, got %s", string(body))
	}
}

func TestResponseWriter_CustomStatus_WithData(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]int{"num": 42}

	ResponseWriter(rr, data, http.StatusCreated)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), `"num":42`) {
		t.Errorf("expected body to contain num:42, got %s", string(body))
	}
}

func TestResponseWriter_WithNilData(t *testing.T) {
	rr := httptest.NewRecorder()

	ResponseWriter(rr, "success", http.StatusAccepted)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Errorf("expected body, got %s", string(body))
	}
}

func TestResponseWriter_MarshalError(t *testing.T) {
	rr := httptest.NewRecorder()
	ch := make(chan int) // not serializable

	ResponseWriter(rr, ch, 0)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "json: unsupported type") {
		t.Errorf("expected marshal error, got %s", string(body))
	}
}
