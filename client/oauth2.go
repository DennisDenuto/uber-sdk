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

func (oauth *Oauth2) Get(url string, queryParams map[string]string) (io.Reader, error) {
	if !oauth.AccessToken.Valid() {
		updatedToken, err := oauth.TokenSource(context.Background(), oauth.AccessToken).Token()
		if err != nil {
			return nil, err
		}

		oauth.AccessToken = updatedToken
	}

	resp, err := oauth.Client(context.Background(), oauth.AccessToken).Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}


