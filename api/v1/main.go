package main

import (
	"os"

	_ "github.com/JonathanGzzBen/nutrity-api/api/v1/docs"
	"github.com/JonathanGzzBen/nutrity-api/api/v1/repository"
	"github.com/JonathanGzzBen/nutrity-api/api/v1/server"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/driver/postgres"
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
	var db *gorm.DB
	var err error
	if os.Getenv("NUTRITY_DB_POSTGRE") == "true" || os.Getenv("NUTRITY_DB_POSTGRE") == "True" || os.Getenv("NUTRITY_DB_POSTGRE") == "1" {
		dbHost := os.Getenv("NUTRITY_DB_HOST")
		dbUser := os.Getenv("NUTRITY_DB_USER")
		dbPassword := os.Getenv("NUTRITY_DB_PASS")
		dbPort := os.Getenv("NUTRITY_DB_PORT")
		if dbHost == "" || dbUser == "" || dbPassword == "" || dbPort == "" {
			panic("Missing database configuration environment variables. See .env.example")
		}
		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " port=" + dbPort + " dbname=nutrity sslmode=disable"
		db, err = gorm.Open(postgres.Open(dsn))
	} else {
		db, err = gorm.Open(sqlite.Open("nutrity.db"))
	}
	if err != nil {
		panic("Could not connect to database")
	}
	serverConfig := server.ServerConfig{
		GoogleConfig: &oauth2.Config{
			ClientID:     os.Getenv("NUTRITY_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("NUTRITY_GOOGLE_CLIENT_SECRET"),
			Endpoint:     endpoints.Google,
			RedirectURL:  "http://127.0.0.1:8080/v1/auth/google-callback",
			Scopes:       []string{"openid", "profile", "email"},
		},
		UsersRepo: repository.NewUsersGormRepository(db),
	}
	// hostname is used by multiple controllers
	// to make requests to authentication controller
	hostname := os.Getenv("NUTRITY_HOSTNAME")
	if hostname == "" {
		panic("Environment variable NUTRITY_HOSTNAME missing")
	}
	serverConfig.Hostname = hostname
	s := server.NewServer(serverConfig)

	port := os.Getenv("NUTRITY_PORT")
	if port == "" {
		s.Run("80")
	} else {
		s.Run(port)
	}
}
