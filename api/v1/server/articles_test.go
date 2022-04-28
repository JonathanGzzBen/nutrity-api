package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"github.com/JonathanGzzBen/ingenialists/api/v1/repository/mocks"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
)

func TestGetAllArticles(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	mockArticles := []models.Article{
		models.Article{
			ID:    1,
			Title: "First article",
		},
		models.Article{
			ID:    2,
			Title: "Second article",
		},
		models.Article{
			ID:    3,
			Title: "Third article",
		},
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetAllArticles").Return(mockArticles, nil)
	s.ArticlesRepo = mockArticlesRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/articles", ts.URL))
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

	var resArticles []models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticles)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(mockArticles) != len(resArticles) {
		t.Fatalf("Expected %v, got %v", len(mockArticles), len(resArticles))
	}
}

func TestGetArticle(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToGet := &models.Article{
		ID:    1,
		Title: "First article",
	}
	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToGet.ID).Return(aToGet, nil)

	s.ArticlesRepo = mockArticlesRepo

	res, err := http.Get(fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToGet.ID))
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

	var resArticle models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resArticle.Title != aToGet.Title {
		t.Fatalf("Expected %v, got %v", aToGet.Title, resArticle.Title)
	}
}

func TestCreateArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToCreate := &models.Article{
		ID:    0,
		Title: "First article",
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	s.ArticlesRepo = mockArticlesRepo

	maJSONBytes, err := json.Marshal(aToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

func TestCreateArticleAsReaderReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToCreate := models.Article{
		ID:    0,
		Title: "First article",
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	s.ArticlesRepo = mockArticlesRepo

	maJSONBytes, err := json.Marshal(aToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Reader")
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

func TestCreateArticleAsWriterReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	c := models.Category{
		ID:   1,
		Name: "First category",
	}

	aToCreate := models.Article{
		ID:         0,
		Title:      "First article",
		CategoryID: c.ID,
		UserID:     1,
	}

	aCreated := aToCreate
	aCreated.ID = 1

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("CreateArticle", &aToCreate).Return(&aCreated, nil)
	s.ArticlesRepo = mockArticlesRepo
	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("GetCategory", aToCreate.CategoryID).Return(&c, nil)
	s.CategoriesRepo = mockCategoriesRepo

	maJSONBytes, err := json.Marshal(aToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
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

	var resArticle models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resArticle.Title != aToCreate.Title {
		t.Fatalf("Expected %v, got %v", aToCreate.Title, resArticle.Title)
	}
}

func TestCreateArticleAsAdministratorReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	c := models.Category{
		ID:   1,
		Name: "First category",
	}

	aToCreate := models.Article{
		ID:         0,
		Title:      "First article",
		CategoryID: c.ID,
		UserID:     1,
	}

	aCreated := aToCreate
	aCreated.ID = 1

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("CreateArticle", &aToCreate).Return(&aCreated, nil)
	s.ArticlesRepo = mockArticlesRepo
	mockCategoriesRepo := &mocks.CategoriesRepository{}
	mockCategoriesRepo.On("GetCategory", aToCreate.CategoryID).Return(&c, nil)
	s.CategoriesRepo = mockCategoriesRepo

	maJSONBytes, err := json.Marshal(aToCreate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/articles", ts.URL), bytes.NewBuffer(maJSONBytes))
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

	var resArticle models.Article
	err = json.NewDecoder(res.Body).Decode(&resArticle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resArticle.Title != aToCreate.Title {
		t.Fatalf("Expected %v, got %v", aToCreate.Title, resArticle.Title)
	}
}

func TestUpdateArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToUpdate := models.Article{
		ID:         0,
		Title:      "First article",
		CategoryID: 1,
		UserID:     1,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

func TestUpdateArticleAsReaderReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToUpdate := models.Article{
		ID:         0,
		Title:      "First article",
		CategoryID: 1,
		UserID:     1,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
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

func TestUpdateArticleAsWriterThatOwnsArticleReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	c := models.Category{
		ID:   1,
		Name: "First category",
	}

	aToUpdate := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: c.ID,
		UserID:     1,
	}

	aUpdated := aToUpdate
	aUpdated.Title = "Article Updated"

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToUpdate.ID).Return(&aToUpdate, nil)
	mockArticlesRepo.On("UpdateArticle", &aToUpdate).Return(&aUpdated, nil)
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
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
}

func TestUpdateArticleAsWriterThatDoesNotOwnArticleReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToUpdate := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	aUpdated := aToUpdate
	aUpdated.Title = "Article Updated"

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToUpdate.ID).Return(&aToUpdate, nil)
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
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

func TestUpdateArticleAsAdministratorThatOwnsArticleReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToUpdate := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     1,
	}

	aUpdated := aToUpdate
	aUpdated.Title = "Article Updated"

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToUpdate.ID).Return(&aToUpdate, nil)
	mockArticlesRepo.On("UpdateArticle", &aToUpdate).Return(&aUpdated, nil)
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
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
}

func TestUpdateArticleAsAdministratorThatDoesNotOwnArticleReturnOk(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToUpdate := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	aUpdated := aToUpdate
	aUpdated.Title = "Article Updated"

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToUpdate.ID).Return(&aToUpdate, nil)
	mockArticlesRepo.On("UpdateArticle", &aToUpdate).Return(&aUpdated, nil)
	s.ArticlesRepo = mockArticlesRepo

	mcJSONBytes, err := json.Marshal(aToUpdate)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToUpdate.ID), bytes.NewBuffer(mcJSONBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Administrator")
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

func TestDeleteArticleAsUnauthenticatedUserReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     1,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
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

	// Verify that mockArticle is still in database
	aInDB, err := s.ArticlesRepo.GetArticle(aToDelete.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if aInDB.Title != aToDelete.Title {
		t.Fatalf("Expected %v, got %v", aToDelete.Title, aInDB.Title)
	}
}

func TestDeleteArticleAsReaderReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Reader")
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

func TestDeleteArticleAsWriterThatDoesNotOwnArticleReturnForbidden(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
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

func TestDeleteArticleAsWriterThatOwnsArticleReturnNoContent(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     1,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	mockArticlesRepo.On("DeleteArticle", aToDelete.ID).Return(nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Add(server.AccessTokenName, "Writer")
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

func TestDeleteArticleAsAdministratorThatDoesNotOwnArticleReturnNoContent(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	mockArticlesRepo.On("DeleteArticle", aToDelete.ID).Return(nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
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

func TestDeleteArticleAsAdministratorThatOwnsArticleReturnNoContent(t *testing.T) {
	s := NewTestServer()
	ts := httptest.NewServer(s.Router)
	defer ts.Close()

	aToDelete := models.Article{
		ID:         1,
		Title:      "First article",
		CategoryID: 1,
		UserID:     2,
	}

	mockArticlesRepo := &mocks.ArticlesRepository{}
	mockArticlesRepo.On("GetArticle", aToDelete.ID).Return(&aToDelete, nil)
	mockArticlesRepo.On("DeleteArticle", aToDelete.ID).Return(nil)
	s.ArticlesRepo = mockArticlesRepo

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/articles/%d", ts.URL, aToDelete.ID), nil)
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
