package client_test

import (
	. "github.com/DennisDenuto/uber-client/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"golang.org/x/oauth2"
	"time"
	"net/http"
	"github.com/golang/oauth2/uber"
)

var _ = Describe("Oauth2", func() {

	Context("When correct client/secret is set", func() {
		var oauth2Config Oauth2

		BeforeEach(func() {
			oauth2Config = Oauth2{
				oauth2.Config{
					ClientID: "valid-client-id",
					ClientSecret: "valid-client-secret",
					RedirectURL: "https://valid-redirect-url",
					Endpoint: uber.Endpoint,
					Scopes: []string{"profile"},
				},
				nil, nil,
			}
		})

		Describe("AuthorisationTokenUrl", func() {
			Context("When no authorisation code has been set", func() {
				It("Should generate a correct authorisation code url", func() {
					authCodeUrl := oauth2Config.AuthorisationTokenUrl()

					Expect(authCodeUrl).ToNot(BeNil())
					Expect(authCodeUrl).To(Equal("https://login.uber.com/oauth/v2/authorize?access_type=offline&client_id=valid-client-id&redirect_uri=https%3A%2F%2Fvalid-redirect-url&response_type=code&scope=profile&state=state"))
				})
			})
		})

		Describe("Token", func() {
			Context("When an authorisation code has been retrieved", func() {
				var server *ghttp.Server

				BeforeEach(func() {
					server = ghttp.NewServer()
					oauth2Config.Endpoint = oauth2.Endpoint{TokenURL: server.URL()}

					server.AppendHandlers(
						func(w http.ResponseWriter, r *http.Request) {
							ghttp.VerifyRequest("POST", "/")(w, r)
							ghttp.VerifyFormKV("client_id", "valid-client-id")(w, r)

							ghttp.RespondWithJSONEncoded(200, struct {
								AccessToken  string `json:"access_token"`
								RefreshToken string `json:"refresh_token"`
								ExpiresIn    int `json:"expires_in"`
								Scope        string `json:"scope"`
							}{
								"ACCESS_TOKEN",
								"REFRESH_TOKEN",
								60,
								"profile history",
							})(w, r)
						},
					)
				})

				AfterEach(func() {
					server.Close()
				})

				It("Should generate a correct access token", func() {
					token, err := oauth2Config.Token("AUTH_TOKEN_GENERATED")

					Expect(err).ToNot(HaveOccurred())
					Expect(token).ToNot(BeNil())
					Expect(token.AccessToken).To(Equal("ACCESS_TOKEN"))
					Expect(token.RefreshToken).To(Equal("REFRESH_TOKEN"))
					Expect(token.Expiry).Should(BeTemporally("~", time.Now(), time.Minute))

					Expect(oauth2Config.AccessToken.AccessToken).To(Equal("ACCESS_TOKEN"))
					Expect(oauth2Config.AccessToken.RefreshToken).To(Equal("REFRESH_TOKEN"))
				})
			})

		})
	})

	Context("When invalid client/secret is set", func() {
		var server *ghttp.Server

		var oauth2Config Oauth2

		BeforeEach(func() {
			oauth2Config = Oauth2{
				oauth2.Config{
					ClientID: "invalid-client-id",
					ClientSecret: "invalid-client-secret",
					RedirectURL: "https://valid-redirect-url",
					Scopes: []string{"profile"},
					Endpoint: uber.Endpoint,
				},
				nil, nil,
			}
		})

		Describe("Token", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				oauth2Config.Endpoint = oauth2.Endpoint{TokenURL: server.URL()}

				server.AppendHandlers(
					func(w http.ResponseWriter, r *http.Request) {
						ghttp.VerifyRequest("POST", "/")(w, r)

						ghttp.RespondWithJSONEncoded(401, struct {
							Error  string `json:"error"`
						}{
							"invalid_client",
						})(w, r)
					},
				)
			})

			AfterEach(func() {
				server.Close()
			})

			It("Should not set accesstoken and instead return an error", func() {
				_, err := oauth2Config.Token("AUTH_TOKEN_GENERATED")

				Expect(err).To(HaveOccurred())
				Expect(oauth2Config.AccessToken).To(BeNil())
			})
		})
	})

})
