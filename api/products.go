package api

import (
	"github.com/DennisDenuto/uber-sdk/client"
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
	"io/ioutil"
)

type ProductsInfo interface {
	GetProduct(productId string) (Product, error)
}

type Products struct {
	ServerTokenClient client.Client
}

func NewProducts(serverToken string) Products {
	return Products{
		ServerTokenClient: client.ServerTokenClient{
			serverToken: serverToken,
		},
	}
}

func (p Products) GetProduct(productId string) (Product, error) {
	resp, err := p.ServerTokenClient.Get(fmt.Sprintf("v1/products/%s", productId), nil)

	if err != nil {
		return Product{}, errors.Wrap(err, "Unable to get Product Info from uber")
	}

	uberProduct := Product{}
	uberResp, _ := ioutil.ReadAll(resp)
	err = json.Unmarshal(uberResp, &uberProduct)

	if err != nil {
		return Product{}, errors.Wrap(err, "Unable to parse Product response from uber")
	}

	return uberProduct, nil
}

type Product struct {
	ProductID    string `json:"product_id"`

	Description  string `json:"description"`

	DisplayName  string `json:"display_name"`

	Capacity     int `json:"capacity"`

	Image        string `json:"image"`

	Shared       bool `json:"shared"`

	PriceDetails PriceDetails `json:"price_details"`
}

type PriceDetails struct {
	DistanceUnit    string `json:"distance_unit"`

	CostPerMin      float64 `json:"cost_per_minute"`

	ServiceFees     []float64 `json:"service_fees"`

	Minimum         float64 `json:"minimum"`

	CostPerDistance float64 `json:"cost_per_distance"`

	Base            float64 `json:"base"`

	CancellationFee float64 `json:"cancellation_fee"`

	CurrencyCode    string `json:"currency_code"`
}
