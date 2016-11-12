package api_test

import (
	. "github.com/DennisDenuto/uber-client/api"

	"github.com/DennisDenuto/uber-client/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

var _ = Describe("Riders", func() {

	Context("User has valid access token", func() {
		var ridersClient Riders
		var server *ghttp.Server
		var oauth2Config client.Oauth2

		BeforeEach(func() {
			server = ghttp.NewServer()

			oauth2Config.AccessToken = &oauth2.Token{Expiry: time.Now().Add(time.Minute), RefreshToken: "Valid-refresh-token", AccessToken: "Valid-Access-Token"}
			oauth2Config.Endpoint = oauth2.Endpoint{
				TokenURL: server.URL(),
			}
			oauth2Config.RootUrl = server.URL() + "/"

			ridersClient = RiderInfo{
				Oauth2: oauth2Config,
			}
		})

		Describe("Me", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					func(w http.ResponseWriter, r *http.Request) {
						ghttp.VerifyRequest("GET", "/me")(w, r)
						ghttp.RespondWith(200, `{
									  "first_name": "Uber",
									  "last_name": "Developer",
									  "email": "developer@uber.com",
									  "picture": "https://...",
									  "promo_code": "teypo",
									  "mobile_verified": true,
									  "uuid": "91d81273-45c2-4b57-8124-d0165f8240c0"
									}`)(w, r)
					},
				)
			})

			It("Should parse response from uber correctly", func() {
				user, err := ridersClient.Me()
				Expect(err).ToNot(HaveOccurred())

				Expect(user.FirstName).To(Equal("Uber"))
				Expect(user.UUID).To(Equal("91d81273-45c2-4b57-8124-d0165f8240c0"))
			})
		})
	})

})
