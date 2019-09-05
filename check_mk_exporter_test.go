package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckMkHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/check_mk", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckMkHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v. Body:\n%s",
			status, http.StatusOK, rr.Body.String())
	}
}
