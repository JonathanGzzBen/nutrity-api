package server_test

import (
	"context"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

type OAuth2ConfigMock struct{}

func (o *OAuth2ConfigMock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   "v1/auth/authorize",
	}

	v := url.Values{}
	v.Set("state", state)

	u.RawQuery = v.Encode()
	return u.String()
}

func (o *OAuth2ConfigMock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "AccessToken",
		Expiry:      time.Now().Add(1 * time.Hour),
	}, nil
}
