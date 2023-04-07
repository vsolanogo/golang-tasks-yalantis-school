package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"task4/handlers"
	"task4/models"
	"fmt"
)

func TestGetProjectsHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/apitrello/project", handlers.GetProjectsHandler)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/apitrello/project", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var projects []models.Project
	json.Unmarshal(writer.Body.Bytes(), &projects)

	if projects[0].ID != 1 {
		t.Errorf("error")
	}
}
