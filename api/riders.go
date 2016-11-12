package api

import (
	"encoding/json"
	"github.com/DennisDenuto/uber-sdk/client"
	"io/ioutil"
)

type Riders interface {
	Me() (User, error)
	History() (UserActivity, error)
}

type RiderInfo struct {
	Oauth2 client.Oauth2
}

func NewRiderInfo(clientId string, clientSecret string, redirectUrl string) RiderInfo {
	return RiderInfo{client.NewOauth2(clientId, clientSecret, []string{"profile", "history"}, redirectUrl)}
}

func (riderInfo RiderInfo) Me() (User, error) {
	response, err := riderInfo.Oauth2.Get("v1/me", nil)

	if err != nil {
		return User{}, err
	}

	user := User{}
	userBytes, _ := ioutil.ReadAll(response)
	err = json.Unmarshal(userBytes, &user)

	if err != nil {
		return User{}, err
	}

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

func (riderInfo RiderInfo) History() (UserActivity, error){
	response, err := riderInfo.Oauth2.Get("v1.2/history", nil)

	if err != nil {
		return UserActivity{}, err
	}

	userActivity := UserActivity{}
	userActivityBytes, _ := ioutil.ReadAll(response)
	err = json.Unmarshal(userActivityBytes, &userActivity)

	if err != nil {
		return UserActivity{}, err
	}
	return userActivity, err
}

type UserActivity struct {
	Offset int `json:"offset"`

	Limit int `json:"limit"`

	Count int `json:"count"`

	History []*Trip `json:"history"`
}

type Trip struct {
	Uuid string `json:"uuid"`

	RequestTime int `json:"request_time"`

	ProductID string `json:"product_id"`

	Status string `json:"status"`

	Distance float64 `json:"distance"`

	StartTime int `json:"start_time"`

	StartLocation *Location `json:"start_location"`

	EndTime int `json:"end_time"`

	EndLocation *Location `json:"end_location"`
}

type Location struct {
	Address string `json:"address,omitempty"`

	Latitude float64 `json:"latitude"`

	Longitude float64 `json:"longitude"`
}