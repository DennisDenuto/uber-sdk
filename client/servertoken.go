package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ServerTokenClient struct {
	RootUrl     string
	ServerToken string
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
