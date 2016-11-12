package client

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/uber"
	"io"
)

type Oauth2 struct {
	oauth2.Config
	AccessToken *oauth2.Token
	AuthToken   *AuthToken
	RootUrl     string
}

type AuthToken struct {
	AuthCode string
}

func NewOauth2(clientId string, clientSecret string, scopes []string, redirectUrl string) Oauth2 {
	oauth2Client := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     uber.Endpoint,
		Scopes:       scopes,
		RedirectURL:  redirectUrl,
	}

	return Oauth2{
		oauth2Client,
		nil, &AuthToken{}, "",
	}

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

	rootUrl := DefaultRootURL

	if oauth.RootUrl != "" {
		rootUrl = oauth.RootUrl
	}

	getUrl := fmt.Sprintf("%s%s", rootUrl, url)
	resp, err := oauth.Client(context.Background(), oauth.AccessToken).Get(getUrl)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
