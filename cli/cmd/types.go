package cmd

type Opts struct {
	EstimatorCmd `command:"estimate"`
	RidersMeCMd `command:"me"`
}

type RidersMeCMd struct {
	OauthOpts
}

type EstimatorCmd struct {
	GetTimeCmd `command:"get-time"`
	GetPriceCmd `command:"get-price"`
}

type GetTimeCmd struct {
	ServerTokenOpts
	StartLatLonOpts
}

type GetPriceCmd struct {
	ServerTokenOpts
	StartLatLonOpts
	StopLatLonOpts
	SeatCount int `long:"seat-count" required:"true"`
}

type ServerTokenOpts struct {
	ServerToken string `long:"servertoken" description:"server token from uber app account"`
}

type OauthOpts struct {
	ClientId     string `long:"client-id" description:"client-id from uber app account" env:"CLIENT_ID" required:"true"`
	ClientSecret string `long:"client-secret" description:"client-secret from uber app account" env:"CLIENT_SECRET" requred:"true"`
}

type StartLatLonOpts struct {
	StartLat string `long:"start-lat" required:"true"`
	StartLon string `long:"start-lon" required:"true"`
}

type StopLatLonOpts struct {
	StopLat string `long:"stop-lat" required:"true"`
	StopLon string `long:"stop-lon" required:"true"`
}
