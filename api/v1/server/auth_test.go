package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentUser(t *testing.T) {
	e := NewTestEnvironment()
	defer e.Close()
	ts := httptest.NewServer(e.Server.Router)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/v1/auth", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if res.StatusCode != 403 {
		t.Fatalf("Expected status code 403, got %v", res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]

	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-7\", got %s", val[0])
	}
}
