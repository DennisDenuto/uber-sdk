package api

import "github.com/DennisDenuto/uber-sdk/client"

type ProductsInfo interface {
	GetProduct(productId string)
}



type Products struct {
	ServerTokenClient client.Client
}