package server_test

import (
	"os"

	"github.com/JonathanGzzBen/ingenialists/api/v1/repository"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestEnvironment struct {
	Server *server.Server
}

func (e *TestEnvironment) Close() {
	os.Remove("test.db")
}

func NewTestServer() *server.Server {
	os.Remove("test.db")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	return server.NewServer(
		server.ServerConfig{
			GoogleConfig:   &OAuth2ConfigMock{},
			Hostname:       "http://localhost:8080",
			Development:    true,
			CategoriesRepo: repository.NewCategoriesGormRepository(db),
			UsersRepo:      repository.NewUsersGormRepository(db),
			ArticlesRepo:   repository.NewArticlesGormRepository(db),
		},
	)
}

func NewTestEnvironment() *TestEnvironment {
	os.Remove("test.db")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}
	server := server.NewServer(
		server.ServerConfig{
			GoogleConfig:   &OAuth2ConfigMock{},
			Hostname:       "http://localhost:8080",
			Development:    true,
			CategoriesRepo: repository.NewCategoriesGormRepository(db),
			UsersRepo:      repository.NewUsersGormRepository(db),
			ArticlesRepo:   repository.NewArticlesGormRepository(db),
		},
	)
	ts := &TestEnvironment{
		Server: server,
	}
	return ts
}
