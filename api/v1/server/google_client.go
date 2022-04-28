package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

type googleUserInfoResponse struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

type IGoogleClient interface {
	userInfoByAccessToken(string) (*googleUserInfoResponse, error)
}
type GoogleClient struct{}
type GoogleClientMock struct{}

// userInfoByAccessToken returns userInfo
func (g *GoogleClient) userInfoByAccessToken(at string) (*googleUserInfoResponse, error) {
	response, err := http.Get(googleUserInfoURL + "?access_token=" + at)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("invalid access token")
	}
	defer response.Body.Close()
	var uinfo *googleUserInfoResponse
	json.NewDecoder(response.Body).Decode(&uinfo)
	return uinfo, nil
}

// userInfoByAccessToken returns userInfo
func (g *GoogleClientMock) userInfoByAccessToken(at string) (*googleUserInfoResponse, error) {
	switch at {
	case "AccessToken", "Administrator", "Writer", "Reader":
		return &googleUserInfoResponse{
			Sub:  "123123213",
			Name: "Mock User",
		}, nil
	default:
		return nil, errors.New("invalid access token")
	}
}
