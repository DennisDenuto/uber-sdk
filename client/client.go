package client

import "io"

const DefaultRootURL string = "https://api.uber.com/"

type Client interface {
	Get(url string, queryParams map[string]string) (io.Reader, error)
}
