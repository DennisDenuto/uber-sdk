package api_test

import (
	. "github.com/DennisDenuto/uber-sdk/api"

	"github.com/DennisDenuto/uber-sdk/client"
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
						ghttp.VerifyRequest("GET", "/v1/me")(w, r)
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

		Describe("Me", func() {
			var happyJsonResponse string

			BeforeEach(func() {
				happyJsonResponse = `{
							  "offset": 0,
							  "limit": 1,
							  "count": 5,
							  "history": [
							    {
							     "status":"completed",
							     "distance":1.64691465,
							     "request_time":1428876188,
							     "start_time":1428876374,
							     "start_city":{
								"display_name":"San Francisco",
								"latitude":37.7749295,
								"longitude":-122.4194155
							     },
							     "end_time":1428876927,
							     "request_id":"37d57a99-2647-4114-9dd2-c43bccf4c30b",
							     "product_id":"a1111c8c-c720-46c3-8534-2fcdd730040d"
							  }
							  ]
							}`
				server.AppendHandlers(
					func(w http.ResponseWriter, r *http.Request) {
						ghttp.VerifyRequest("GET", "/v1.2/history")(w, r)
						ghttp.RespondWith(200, happyJsonResponse)(w, r)
					},
				)
			})

			It("Should parse response from uber correctly", func() {
				userActivity, err := ridersClient.History()
				Expect(err).ToNot(HaveOccurred())

				Expect(userActivity.Offset).To(Equal(0))
				Expect(userActivity.History).To(HaveLen(1))
				Expect(userActivity.History[0].Status).To(Equal("completed"))
			})

			Context("Malformed json response from uber", func() {

				BeforeEach(func() {
					happyJsonResponse = "INVALID JSON"
				})

				It("should return error", func() {
					_, err := ridersClient.History()

					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

})
