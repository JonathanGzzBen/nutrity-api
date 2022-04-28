package main

import (
	"os"

	_ "github.com/JonathanGzzBen/ingenialists/api/v1/docs"
	"github.com/JonathanGzzBen/ingenialists/api/v1/repository"
	"github.com/JonathanGzzBen/ingenialists/api/v1/server"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Ingenialists API V1
// @version v1.0.0
// @description This is Ingenialist's API
//
// @contact.name JonathanGzzBen
// @contact.url http://www.github.com/JonathanGzzBen
// @contact.email jonathangzzben@gmail.com
// @license.name MIT License
// @license.url https://mit-license.org/
//
// @host localhost:8080
// @BasePath /v1
//
// @securityDefinitions.apikey AccessToken
// @in header
// @name AccessToken
//
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl /v1/auth/google-callback
// @authorizationUrl /v1/auth/google-login
// @scope.openid Allow identifying account
// @scope.profile Grant access to profile
// @scope.email Grant access to email
func main() {
	godotenv.Load(".env")
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic("Could not connect to database")
	}
	serverConfig := server.ServerConfig{
		GoogleConfig: &oauth2.Config{
			ClientID:     os.Getenv("ING_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("ING_GOOGLE_CLIENT_SECRET"),
			Endpoint:     endpoints.Google,
			RedirectURL:  "http://127.0.0.1:8080/v1/auth/google-callback",
			Scopes:       []string{"openid", "profile", "email"},
		},
		CategoriesRepo: repository.NewCategoriesGormRepository(db),
		UsersRepo:      repository.NewUsersGormRepository(db),
		ArticlesRepo:   repository.NewArticlesGormRepository(db),
	}
	// hostname is used by multiple controllers
	// to make requests to authentication controller
	hostname := os.Getenv("ING_HOSTNAME")
	if len(hostname) == 0 {
		panic("Environment variable ING_HOSTNAME missing")
	}
	serverConfig.Hostname = hostname
	s := server.NewServer(serverConfig)

	if os.Getenv("ING_ENVIRONMENT") == "development" {
		port := os.Getenv("ING_PORT")
		if len(port) == 0 {
			panic("Environment variable ING_PORT missing")
		}
		s.Run(port)
	} else {
		s.Run()
	}
}
