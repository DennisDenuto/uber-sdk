package client

import "io"

type Client interface {
	Get(url string, queryParams map[string]string) (io.Reader, error)
}