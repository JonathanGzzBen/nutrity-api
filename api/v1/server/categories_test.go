package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/repository/mocks"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
)

// This data should not be modified, its purpose
// is to be used to initialize database.
var mockCategories = []models.Category{
	{ID: 1231, Name: "First Category", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
	{ID: 2131, Name: "Second Name", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
	{ID: 56232, Name: "Third Name", ImageURL: "https://i.imgur.com/oCsJWt7.jpeg"},
}

func TestGetAllCategories(t *testing.T) {
	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("GetAllCategories").Return(mockCategories, nil)
	s := NewTestServer()
	s.CategoriesRepo = mockCategoriesRepo
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/v1/categories", ts.URL))
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

	var resCategories []models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategories)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockCategories) != len(resCategories) {
		t.Fatalf("Expected %v, got %v", len(mockCategories), len(resCategories))
	}
}

func TestGetCategory(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockCategory := &mockCategories[1]
	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("GetCategory", mockCategory.ID).Return(mockCategory, nil)

	s.CategoriesRepo = mockCategoriesRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/categories/"+strconv.Itoa(int(mockCategory.ID)), ts.URL))
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
	var resCategory *models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(resCategory, mockCategory) {
		t.Fatalf("Expected %v, got %v", mockCategory, resCategory)
	}

}

func TestCreateCategoryAsRegularUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToCreate := mockCategories[1]

	mockCategoriesRepo := &mocks.CategoriesRepository{}
	s.CategoriesRepo = mockCategoriesRepo

	mcJSONBytes, err := json.Marshal(cToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/categories", ts.URL), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "AccessToken")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestCreateCategoryAsWriterReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToCreate := mockCategories[1]
	cToCreate.ID = 0

	mockCategoriesRepo := &mocks.CategoriesRepository{}
	s.CategoriesRepo = mockCategoriesRepo

	mcJSONBytes, err := json.Marshal(cToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/categories", ts.URL), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestCreateCategoryAsAdministratorReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToCreate := mockCategories[1]
	cToCreate.ID = 0

	cCreated := cToCreate
	cCreated.ID = 1

	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("CreateCategory", &cToCreate).Return(&cCreated, nil)
	s.CategoriesRepo = mockCategoriesRepo

	mcJSONBytes, err := json.Marshal(cToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/categories", ts.URL), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

	var resCategory *models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(*resCategory, cCreated) {
		t.Fatalf("Expected %v, got %v", cCreated, *resCategory)
	}
}

func TestUpdateCategoryAsRegularUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToUpdate := mockCategories[1]
	cToUpdate.Name = "Category Updated"

	mcJSONBytes, err := json.Marshal(cToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/categories/%d", ts.URL, cToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "AccessToken")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestUpdateCategoryAsWriterReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToUpdate := mockCategories[1]
	cToUpdate.Name = "Category Updated"

	mcJSONBytes, err := json.Marshal(cToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/categories/%d", ts.URL, cToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

}

func TestUpdateCategoryAsAdministratorReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	cToUpdate := mockCategories[1]
	cToUpdate.Name = "Category Updated"

	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("UpdateCategory", &cToUpdate).Return(&cToUpdate, nil)
	mockCategoriesRepo.On("GetCategory", cToUpdate.ID).Return(&cToUpdate, nil)
	s.CategoriesRepo = mockCategoriesRepo

	mcJSONBytes, err := json.Marshal(cToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/categories/%d", ts.URL, cToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

	var resCategory *models.Category
	err = json.NewDecoder(res.Body).Decode(&resCategory)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(*resCategory, cToUpdate) {
		t.Fatalf("Expected %v, got %v", cToUpdate, *resCategory)
	}
}

func TestDeleteCategoryAsRegularUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockCategory := mockCategories[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/categories/%d", ts.URL, mockCategory.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Access Token")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}
}

func TestDeleteCategoryAsWriterReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockCategory := mockCategories[1]

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/categories/%d", ts.URL, mockCategory.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code %v, got %v", http.StatusForbidden, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "application/json; charset=utf-8" {
		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	}

}

func TestDeleteCategoryAsAdministratorReturnNoContent(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockCategory := mockCategories[1]

	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("DeleteCategory", mockCategory.ID).Return(nil)
	mockCategoriesRepo.On("GetCategory", mockCategory.ID).Return(&mockCategory, nil)
	s.CategoriesRepo = mockCategoriesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/categories/%d", ts.URL, mockCategory.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, res.StatusCode)
	}

	val, ok := res.Header["Content-Type"]
	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}
	if val[0] != "text/plain; charset=utf-8" {
		t.Fatalf("Expected \"text/plain; charset=utf-8\", got %s", val[0])
	}
}
