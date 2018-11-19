package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/maguec/micro-leaderboard/handlers/app"
	//	"github.com/maguec/micro-leaderboard/handlers/healthcheck"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	type ApiMessage struct {
		Message string
	}
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/", app.Root)
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("HTTP response code should be 200, was: %d", w.Code)
	}
	var message ApiMessage
	err := json.Unmarshal(w.Body.Bytes(), &message)
	if err != nil {
		t.Errorf("JSON parse error: %d", w.Code)
	}

	if message.Message != "This is root - Please see the docs" {
		t.Errorf("JSON message should be OK, but is: %s", message.Message)
	}
}

// TODO: figure out a way to pass the redis client information
//func TestHealth(t *testing.T) {
//	w := httptest.NewRecorder()
//	r := gin.Default()
//	r.GET("/health", healthcheck.HealthCheck)
//	req, _ := http.NewRequest("GET", "/health", nil)
//	r.ServeHTTP(w, req)
//	if w.Code != 500 {
//		t.Errorf("HTTP response code should be 200, was: %d", w.Code)
//	}
//}
