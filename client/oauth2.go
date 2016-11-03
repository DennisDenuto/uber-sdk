package client

import (
	"io"
	"golang.org/x/oauth2"
	"context"
	"time"
)

type Oauth2 struct {
	oauth2.Config
	AccessToken  *accessToken
	AuthToken    *authToken
}

type authToken struct {
	AuthCode string
}

func (oauth *Oauth2) AuthorisationTokenUrl() string {
	return oauth.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

type accessToken struct {
	AccessToken  string
	RefreshToken string
	Expiry time.Time
}

func (oauth *Oauth2) Token(authToken string) (*oauth2.Token, error) {
	token, err := oauth.Exchange(context.Background(), authToken)

	if err != nil {
		return nil, err
	}

	oauth.AccessToken = &accessToken{
		AccessToken: token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry: token.Expiry,
	}

	return token, nil
}

func (oauth2 Oauth2) Get(url string, queryParams map[string]string) (io.Reader, error) {
	return nil, nil
}


