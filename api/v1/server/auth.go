package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/JonathanGzzBen/nutrity-api/api/v1/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	state             = "nutrity-api"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	AccessTokenName   = "AccessToken"
	TokenLength       = 40
)

type IOauthConfig interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

type tokenResponse struct {
	AccessToken string `json:"AccessToken"`
	TokenType   string `json:"TokenType"`
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
	u, err := s.userByAccessToken(at)
	if err != nil {
		c.JSON(http.StatusForbidden, models.APIError{Code: http.StatusForbidden, Message: "invalid access token"})
		return
	}
	c.JSON(http.StatusOK, userDTOFromUser(u))
}

// LoginGoogle is the handler for GET requests to /auth/google-login
// it's the entryway for Google OAuth2 flow.
func (s *Server) LoginGoogle(c *gin.Context) {
	url := s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback is the handler for GET requests to /auth/google-callback
// it's part of Google OAuth2 flow.
//
// Returns user's token.
func (s *Server) GoogleCallback(c *gin.Context) {
	if c.Request.URL.Query().Get("state") != state {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "state did not match"})
		return
	}

	authCode := c.Request.URL.Query().Get("code")
	ctx := context.Background()
	token, err := s.googleConfig.Exchange(ctx, authCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to exchange token: " + err.Error()})
		return
	}

	uinfo, err := s.googleClient.userInfoByAccessToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.APIError{Code: http.StatusBadRequest, Message: "failed to get user info: " + err.Error()})
		return
	}

	u, err := s.UsersRepo.GetUserByGoogleSub(uinfo.Sub)
	tokenResponse := tokenResponse{
		TokenType: "Bearer",
	}
	if u != nil {
		tokenResponse.AccessToken = u.AccessToken
	}
	if err != nil {
		u := &models.User{
			GoogleSub:   uinfo.Sub,
			AccessToken: generateSecureToken(TokenLength),
			Username:    uinfo.Name,
			Email:       uinfo.Email,
		}
		s.UsersRepo.CreateUser(u)
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (s *Server) userByAccessToken(at string) (*models.User, error) {
	u, err := s.UsersRepo.GetUserByAccessToken(at)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// devOAuthAuthorize handles requests to /auth/authorize
// should only be available during development
func (s *Server) devOAuthAuthorize(c *gin.Context) {
	state := c.Request.URL.Query().Get("state")

	u, err := url.Parse("http:://localhost/callback")
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			models.APIError{
				Code:    http.StatusInternalServerError,
				Message: "could not make url",
			})
	}

	v := url.Values{}
	v.Set("code", "code")
	v.Set("state", state)
	u.RawQuery = v.Encode()

	http.Redirect(c.Writer, c.Request, u.String(), http.StatusTemporaryRedirect)
}
