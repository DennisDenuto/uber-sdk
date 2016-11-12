# uber-sdk

## Description

Golang bindings for Uber API.

# Usage

## Register your app
Register an application at the [Uber Developer Portal](https://developer.uber.com). You will receive a client_id, secret, and server_token.

## API requests

All the APIs are available from the `api` package.
```
api
├── estimates.go
├── products.go
├── riders.go
```

```go
estimator := api.NewEstimate("Your-Server-Token")
priceResp, err := estimator.GetTime(c.StartLon, c.StartLat)
fmt.Println(priceResp)
```

# CLI
```
go get github.com/DennisDenuto/uber-sdk

uber-sdk --help
```
