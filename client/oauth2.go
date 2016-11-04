package client

import (
	"io"
	"golang.org/x/oauth2"
	"context"
)

type Oauth2 struct {
	oauth2.Config
	AccessToken *oauth2.Token
	AuthToken   *AuthToken
}

type AuthToken struct {
	AuthCode string
}

func (oauth *Oauth2) AuthorisationTokenUrl() string {
	return oauth.Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (oauth *Oauth2) Token(authToken string) (*oauth2.Token, error) {
	token, err := oauth.Exchange(context.Background(), authToken)

	if err != nil {
		return nil, err
	}

	oauth.AccessToken = token

	return token, nil
}

func (oauth2 Oauth2) Get(url string, queryParams map[string]string) (io.Reader, error) {
	resp, _ := oauth2.Client(context.Background(), oauth2.AccessToken).Get(url)

	return resp.Body, nil
}


