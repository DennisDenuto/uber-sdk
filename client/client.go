package client

import "io"

const DefaultRootURL string = "https://api.uber.com/v1/"


type Client interface {
	Get(url string, queryParams map[string]string) (io.Reader, error)
}