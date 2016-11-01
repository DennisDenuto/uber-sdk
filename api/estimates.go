package api

import (
	"github.com/DennisDenuto/uber-client/client"
	"encoding/json"
	"io/ioutil"
	"github.com/pkg/errors"
)

type Estimator interface {
	GetTime(startLon string, startLat string) (*TimesResp, error)
	GetPrice()
}

type Estimate struct {
	ServerTokenClient client.Client
}

func (estimateClient Estimate) GetTime(startLon string, startLat string) (*TimesResp, error) {
	queryParams := map[string]string{
		"start_latitude" : startLat,
		"start_longitude" : startLon,
	}
	respReader, err := estimateClient.ServerTokenClient.Get("estimates/time", queryParams)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Time Estimates from uber")
	}

	resp, _ := ioutil.ReadAll(respReader)

	var timesResp *TimesResp = &TimesResp{}
	json.Unmarshal(resp, timesResp)

	return timesResp, nil
}

func (estimateClient Estimate) GetPrice() {
	estimateClient.ServerTokenClient.Get("", nil)
}

type TimesResp struct {
	Times []*Time `json:"times"`
}

type Time struct {
	ProductID   string `json:"product_id"`
	DisplayName string `json:"display_name"`
	Estimate    int `json:"estimate"`
}
