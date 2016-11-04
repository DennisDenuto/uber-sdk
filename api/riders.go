package api

import "github.com/DennisDenuto/uber-client/client"

type Riders interface {
	Me()
	History()
}

type RiderInfo struct {
	Oauth2 client.Oauth2
}

func (riderInfo RiderInfo) Me() {
	riderInfo.Oauth2.Get("/v1.2/history", nil)
}

