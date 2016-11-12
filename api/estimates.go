package api

import (
	"encoding/json"
	"github.com/DennisDenuto/uber-client/client"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
)

type Estimator interface {
	GetTime(startLon string, startLat string) (*TimesResp, error)
	GetPrice(startLon string, endLon string, startLat string, endLat string, seatCount int) (*PriceResp, error)
}

type Estimate struct {
	ServerTokenClient client.Client
}

func NewEstimate(serverToken string) Estimate {
	return Estimate{client.ServerTokenClient{ServerToken: serverToken}}
}

func (estimateClient Estimate) GetTime(startLon string, startLat string) (*TimesResp, error) {
	queryParams := map[string]string{
		"start_latitude":  startLat,
		"start_longitude": startLon,
	}
	respReader, err := estimateClient.ServerTokenClient.Get("v1/estimates/time", queryParams)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Time Estimates from uber")
	}

	resp, _ := ioutil.ReadAll(respReader)

	var timesResp *TimesResp = &TimesResp{}
	err = json.Unmarshal(resp, timesResp)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse Time Estimates response from uber")
	}

	return timesResp, nil
}

func (estimateClient Estimate) GetPrice(startLon string, endLon string, startLat string, endLat string, seatCount int) (*PriceResp, error) {
	queryParams := map[string]string{
		"start_latitude":  startLat,
		"start_longitude": startLon,
		"end_latitude":    endLat,
		"end_longitude":   endLon,
		"seat_count":      strconv.Itoa(seatCount),
	}

	respReader, err := estimateClient.ServerTokenClient.Get("v1/estimates/price", queryParams)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Price Estimates from uber")
	}

	resp, _ := ioutil.ReadAll(respReader)

	var priceResp *PriceResp = &PriceResp{}
	err = json.Unmarshal(resp, priceResp)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse Price Estimates response from uber")
	}

	return priceResp, nil
}

type TimesResp struct {
	Times []*Time `json:"times"`
}

type Time struct {
	ProductID   string `json:"product_id"`
	DisplayName string `json:"display_name"`
	Estimate    int    `json:"estimate"`
}

type PriceResp struct {
	Prices []*Price `json:"prices"`
}

type Price struct {
	ProductID       string  `json:"product_id"`
	CurrencyCode    string  `json:"currency_code"`
	DisplayName     string  `json:"display_name"`
	Estimate        string  `json:"estimate"`
	LowEstimate     int     `json:"low_estimate"`
	HighEstimate    int     `json:"high_estimate"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
}
