package api

import (
	"encoding/json"
	"github.com/DennisDenuto/uber-client/client"
	"io/ioutil"
)

type Riders interface {
	Me() (User, error)
	History()
}

type RiderInfo struct {
	Oauth2 client.Oauth2
}

func NewRiderInfo(clientId string, clientSecret string, redirectUrl string) RiderInfo {
	return RiderInfo{client.NewOauth2(clientId, clientSecret, []string{"profile"}, redirectUrl)}
}

func (riderInfo RiderInfo) Me() (User, error) {
	response, err := riderInfo.Oauth2.Get("me", nil)

	if err != nil {
		return User{}, err
	}

	user := User{}
	userBytes, _ := ioutil.ReadAll(response)
	json.Unmarshal(userBytes, &user)
	return user, err
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	PromoCode string `json:"promo_code"`
	UUID      string `json:"uuid"`
}

func (riderInfo RiderInfo) History() {
}
