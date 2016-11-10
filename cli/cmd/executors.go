package cmd

import (
	"fmt"
	"github.com/DennisDenuto/uber-client/api"
	"github.com/DennisDenuto/uber-client/client"
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
	ridersClient := api.RiderInfo{client.NewOauth2(c.ClientId, c.ClientSecret, []string{"profile"}, "https://localhost"), }

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

func (c EstimatorCmd) Execute(args []string) error {
	fmt.Println("ESTIMATOR EXECUTING")
	fmt.Println(c.StartLat)

	return nil
}

