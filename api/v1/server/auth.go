package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/JonathanGzzBen/nutrity-api/api/v1/models"
	"github.com/gin-gonic/gin"
)

var (
	AccessTokenName = "AccessToken"
	TokenLength     = 40
)

type tokenResponse struct {
	AccessToken string `json:"AccessToken"`
	TokenType   string `json:"TokenType"`
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

type RegisterUserDTO struct {
	GoogleToken string `json:"GoogleToken"`
	Username    string `json:"Username"`
	Email       string `json:"Email"`
}

// Register is the handler for POST requests to /auth
// 	@ID RegisterUser
// 	@Tags auth
// 	@Param user body RegisterUserDTO true "User"
// 	@Success 200 {object} tokenResponse
// 	@Failure 403 {object} models.APIError
// 	@Router /auth [post]
func (s *Server) RegisterUser(c *gin.Context) {
	var ru RegisterUserDTO
	if err := c.ShouldBindJSON(&ru); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{Code: http.StatusNotFound, Message: "invalid register user: " + err.Error()})
		return
	}
	fmt.Print("Google token: ")
	fmt.Println(ru.GoogleToken)

	u, err := s.UsersRepo.GetUserByGoogleToken(ru.GoogleToken)
	if err != nil {
		// If no user registered with given GoogleToken
		generatedToken := generateSecureToken(TokenLength)
		u = &models.User{
			AccessToken: generatedToken,
			GoogleToken: ru.GoogleToken,
			Username:    ru.Username,
			Email:       ru.Email,
		}
		s.UsersRepo.CreateUser(u)
	}

	tokenResponse := tokenResponse{
		AccessToken: u.AccessToken,
		TokenType:   "Bearer",
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// CurrentUser is the handler for GET requests to /auth
// 	@ID GetCurrentUser
// 	@Tags auth
// 	@Success 200 {object} UserDTO
// 	@Failure 403 {object} models.APIError
// 	@Security AccessToken
// 	@Router /auth [get]
func (s *Server) GetCurrentUser(c *gin.Context) {
	at := c.GetHeader(AccessTokenName)
	u, err := s.UsersRepo.GetUserByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "invalid access token"})
		return
	}
	c.JSON(http.StatusOK, userDTOFromUser(u))
}
