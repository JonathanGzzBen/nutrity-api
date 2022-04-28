package server

import (
	"github.com/JonathanGzzBen/nutrity-api/api/v1/repository"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Server struct {
	googleClient IGoogleClient
	googleConfig IOauthConfig
	development  bool
	Router       *gin.Engine
	UsersRepo    repository.UsersRepository
}

type ServerConfig struct {
	GoogleConfig IOauthConfig
	Hostname     string
	Development  bool
	UsersRepo    repository.UsersRepository
}

func NewServer(sc ServerConfig) *Server {
	server := &Server{
		googleConfig: sc.GoogleConfig,
		development:  sc.Development,
		UsersRepo:    sc.UsersRepo,
	}
	if sc.Development {
		server.googleClient = &GoogleClientMock{}
	} else {
		server.googleClient = &GoogleClient{}
	}

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		ur := v1.Group("/users")
		{
			ur.GET("/", server.GetAllUsers)
			ur.GET("/:id", server.GetUser)
			ur.PUT("/:id", server.UpdateUser)
		}
		ar := v1.Group("/auth")
		{
			ar.GET("/", server.GetCurrentUser)
			ar.GET("/google-login", server.LoginGoogle)
			ar.GET("/google-callback", server.GoogleCallback)
			if sc.Development {
				ar.GET("/dev-authorize", server.devOAuthAuthorize)
			}
		}
	}

	swaggerUrl := ginSwagger.URL(sc.Hostname + "/v1/swagger/doc.json")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))
	server.Router = router
	return server
}

func (s *Server) Run(port ...string) {
	s.Router.Run(port[0])
}
