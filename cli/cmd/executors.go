package cmd

import (
	"fmt"
	"github.com/DennisDenuto/uber-client/api"
	log "github.com/Sirupsen/logrus"
	"strings"
	"os"
	"bufio"
	"github.com/gosuri/uitable"
)


func AskForAuthCode(authCodeUrl string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(authCodeUrl)
	fmt.Print("\nEnter Authorisation Code: ")

	text, _ := reader.ReadString('\n')

	return strings.TrimSuffix(text, "\n")
}

func (c RidersMeCMd) Execute(args []string) error {
	ridersClient := api.NewRiderInfo(c.ClientId, c.ClientSecret, "https://localhost")

	authCode := AskForAuthCode(ridersClient.Oauth2.AuthorisationTokenUrl())

	_, err := ridersClient.Oauth2.Token(authCode)
	if err != nil {
		log.WithError(err).Error("Could not get an access token from uber")
		return err
	}

	user, err := ridersClient.Me()
	if err != nil {
		log.WithError(err).Error("Could not get profile info from uber")
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("First Name", "Last Name", "Email", "Picture", "Promocode", "UUID")
	table.AddRow(user.FirstName, user.LastName, user.Email, user.Picture, user.PromoCode, user.UUID)
	fmt.Println(table)

	return nil
}

func (c GetTimeCmd) Execute(args []string) error {
	estimator := api.NewEstimate(c.ServerToken)

	priceResp, err := estimator.GetTime(c.StartLon, c.StartLat)
	if err != nil {
		log.WithError(err).Error("Could not estimate time from uber")
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("Display Name", "Estimate", "Product Id")

	for _, value := range priceResp.Times {
		table.AddRow(value.DisplayName, value.Estimate, value.ProductID)
	}
	fmt.Println(table)

	return nil
}

func (c GetPriceCmd) Execute(args []string) error {
	estimator := api.NewEstimate(c.ServerToken)

	priceResp, err := estimator.GetPrice(c.StartLon, c.StartLat, c.StopLon, c.StopLat, c.SeatCount)
	if err != nil {
		log.WithError(err).Error("Could not estimate price from uber")
		return err
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("Display Name", "Product Id", "Currency Code", "Estimate", "Low Estimate", "High Estimate", "Surge Multiply")

	for _, value := range priceResp.Prices {
		table.AddRow(
			value.DisplayName,
			value.ProductID,
			value.CurrencyCode,
			value.Estimate,
			value.LowEstimate,
			value.HighEstimate,
			value.SurgeMultiplier,
		)
	}
	fmt.Println(table)

	return nil
}

