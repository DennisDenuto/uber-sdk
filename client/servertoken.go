package client

import (
	"net/http"
	"io"
	"net/url"
	"fmt"
)

const DefaultRootURL string = "https://api.uber.com/v1/"

type ServerTokenClient struct {
	RootUrl     string
	ServerToken string
	HttpClient  http.Client
}

func (serverTokenClient ServerTokenClient) Get(uberResourceUrl string, queryParams map[string]string) (io.Reader, error) {
	rootUrl := DefaultRootURL

	if serverTokenClient.RootUrl != "" {
		rootUrl = serverTokenClient.RootUrl
	}

	values := url.Values{}
	values.Add("server_token", serverTokenClient.ServerToken)

	for key, param := range queryParams {
		values.Add(key, param)
	}

	getUrl := fmt.Sprintf("%s%s?%s", rootUrl, uberResourceUrl, values.Encode())
	resp, err := http.DefaultClient.Get(getUrl)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}