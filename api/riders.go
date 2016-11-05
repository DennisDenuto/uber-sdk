package api

import (
	"github.com/DennisDenuto/uber-client/client"
	"encoding/json"
	"io/ioutil"
)

type Riders interface {
	Me() (User, error)
	History()
}

type RiderInfo struct {
	Oauth2 client.Oauth2
}

func (riderInfo RiderInfo) Me() (User, error) {
	response, err := riderInfo.Oauth2.Get("/v1.2/me", nil)

	user := User{}
	userBytes, _ := ioutil.ReadAll(response)
	json.Unmarshal(userBytes, &user)
	return user, err
}

type User struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Picture string `json:"picture"`
	PromoCode string `json:"promo_code"`
	UUID string `json:"uuid"`
}

func (riderInfo RiderInfo) History() {
}

