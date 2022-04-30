package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonathanGzzBen/nutrity-api/api/v1/models"
	"github.com/JonathanGzzBen/nutrity-api/api/v1/repository/mocks"
	"github.com/JonathanGzzBen/nutrity-api/api/v1/server"
)

var mockUsers = []models.User{
	{ID: 123, Username: "First User"},
	{ID: 456, Username: "Second User"},
	{ID: 789, Username: "Third User"},
}

func TestGetAllUsers(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	for _, u := range mockUsers {
		s.UsersRepo.CreateUser(&u)
	}

	mockUsersRepo := &mocks.MockUsersRepository{}
	mockUsersRepo.EXPECT().GetAllUsers().Return(mockUsers, nil)
	s.UsersRepo = mockUsersRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/users", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

	var resUsers []models.User
	err = json.NewDecoder(res.Body).Decode(&resUsers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockUsers) != len(resUsers) {
		t.Fatalf("Expected %v, got %v", len(mockUsers), len(resUsers))
	}
}

func TestGetUser(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	uToGet := mockUsers[1]

	mockUsersRepo := &mocks.MockUsersRepository{}
	mockUsersRepo.EXPECT().GetUser(uToGet.ID).Return(&uToGet, nil)
	s.UsersRepo = mockUsersRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/users/%d", ts.URL, uToGet.ID))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

	var resUser models.User
	err = json.NewDecoder(res.Body).Decode(&resUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if uToGet != resUser {
		t.Fatalf("Expected %v, got %v", uToGet, resUser)
	}
}

// TestUpdateUserChangeNameAsDifferentUserReturnForbidden tests a request
// in which a user with a different role from Administrator tries to update
// a user with a different ID than his own.
//
// In testing mode, authenticated user will have ID = 1.
func TestUpdateUserChangeNameAsDifferentUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	// uToUpdate has ID different from 1
	uToUpdate := mockUsers[1]
	uUpdated := uToUpdate
	uUpdated.Username = "Updated name"

	s.UsersRepo = &mocks.MockUsersRepository{}

	muJSONBytes, err := json.Marshal(uUpdated)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, uToUpdate.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %d, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

// TestUpdateUserChangeNameAsSameUserReturnOk tests a request
// in which a user with a different role from Administrator tries to update
// a user with a different ID than his own.
//
// In testing mode, authenticated user will have ID = 1.
func TestUpdateUserChangeNameAsSameUserReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	uToUpdate := mockUsers[1]
	uToUpdate.ID = 1
	uUpdated := uToUpdate
	uUpdated.Username = "Updated username"

	mockUsersRepo := &mocks.MockUsersRepository{}
	mockUsersRepo.EXPECT().GetUser(uToUpdate.ID).Return(&uToUpdate, nil)
	mockUsersRepo.EXPECT().UpdateUser(&uUpdated).Return(&uUpdated, nil)
	s.UsersRepo = mockUsersRepo

	muJSONBytes, err := json.Marshal(uUpdated)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/users/%d", ts.URL, uToUpdate.ID), bytes.NewBuffer(muJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "AccessToken")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %v", http.StatusOK, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

	var resUser models.User
	err = json.NewDecoder(res.Body).Decode(&resUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resUser.Username != uToUpdate.Username {
		t.Fatalf("Expected %v, got %v", uToUpdate.Username, resUser.Username)
	}
}
