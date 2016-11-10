package cmd

type Opts struct {
	EstimatorCmd `command:"estimate"`
	RidersMeCMd `command:"me"`
}

type RidersMeCMd struct {
	OauthOpts
}

type EstimatorCmd struct {
	ServerTokenOpts
	StartLatLonOpts
}

type ServerTokenOpts struct {
	ServerToken string `long:"servertoken" description:"server token from uber app account"`
}

type OauthOpts struct {
	ClientId     string `long:"client-id" description:"client-id from uber app account" env:"CLIENT_ID" required:"true"`
	ClientSecret string `long:"client-secret" description:"client-secret from uber app account" env:"CLIENT_SECRET" requred:"true"`
}

type StartLatLonOpts struct {
	StartLat string `long:"start-lat"`
	StartLon string `long:"start-lon"`
}
