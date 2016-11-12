package client_test

import (
	. "github.com/DennisDenuto/uber-sdk/client"

	"github.com/golang/oauth2/uber"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
)

var _ = Describe("Oauth2", func() {

	Context("When correct client/secret is set", func() {
		var oauth2Config Oauth2

		BeforeEach(func() {
			oauth2Config = Oauth2{
				oauth2.Config{
					ClientID:     "valid-client-id",
					ClientSecret: "valid-client-secret",
					RedirectURL:  "https://valid-redirect-url",
					Endpoint:     uber.Endpoint,
					Scopes:       []string{"profile"},
				},
				nil, nil, "",
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
					oauth2Config.RootUrl = server.URL()

					server.AppendHandlers(
						func(w http.ResponseWriter, r *http.Request) {
							ghttp.VerifyRequest("POST", "/")(w, r)
							ghttp.VerifyFormKV("client_id", "valid-client-id")(w, r)

							ghttp.RespondWithJSONEncoded(200, struct {
								AccessToken  string `json:"access_token"`
								RefreshToken string `json:"refresh_token"`
								ExpiresIn    int    `json:"expires_in"`
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

		Describe("Get", func() {

			Context("access token has been retrieved and valid", func() {
				var server *ghttp.Server

				BeforeEach(func() {
					server = ghttp.NewServer()

					oauth2Config.AccessToken = &oauth2.Token{AccessToken: "Valid-Access-Token"}
					oauth2Config.RootUrl = server.URL()

					server.AppendHandlers(
						func(w http.ResponseWriter, r *http.Request) {
							ghttp.VerifyHeaderKV("Authorization", "Bearer Valid-Access-Token")(w, r)
							ghttp.VerifyRequest("GET", "/")(w, r)
						},
					)
				})

				AfterEach(func() {
					server.Close()
				})

				It("Should use access token when performing Get Request", func() {
					oauth2Config.Get("/", nil)
					Expect(server.ReceivedRequests()).To(HaveLen(1))
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
					ClientID:     "invalid-client-id",
					ClientSecret: "invalid-client-secret",
					RedirectURL:  "https://valid-redirect-url",
					Scopes:       []string{"profile"},
					Endpoint:     uber.Endpoint,
				},
				nil, nil, "",
			}
		})

		Describe("Token", func() {
			BeforeEach(func() {
				server = ghttp.NewServer()
				oauth2Config.Endpoint = oauth2.Endpoint{TokenURL: server.URL()}
				oauth2Config.RootUrl = server.URL()

				server.AppendHandlers(
					func(w http.ResponseWriter, r *http.Request) {
						ghttp.VerifyRequest("POST", "/")(w, r)

						ghttp.RespondWithJSONEncoded(401, struct {
							Error string `json:"error"`
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

		Describe("Get", func() {
			Context("access token has expired and refresh token is present", func() {
				var server *ghttp.Server

				BeforeEach(func() {
					server = ghttp.NewServer()

					oauth2Config.AccessToken = &oauth2.Token{Expiry: time.Now().Add(-time.Minute), RefreshToken: "Valid-refresh-token", AccessToken: "InValid-Access-Token"}
					oauth2Config.AuthToken = &AuthToken{AuthCode: "AuthCode"}
					oauth2Config.Endpoint = oauth2.Endpoint{
						TokenURL: server.URL(),
					}
					oauth2Config.RootUrl = server.URL()

					server.AppendHandlers(
						func(w http.ResponseWriter, r *http.Request) {
							ghttp.VerifyRequest("POST", "/")(w, r)
							ghttp.VerifyFormKV("grant_type", "refresh_token")(w, r)

							ghttp.RespondWithJSONEncoded(200, struct {
								AccessToken  string `json:"access_token"`
								RefreshToken string `json:"refresh_token"`
								ExpiresIn    int    `json:"expires_in"`
								Scope        string `json:"scope"`
							}{
								"NEW_VALID_ACCESS_TOKEN",
								"REFRESH_TOKEN",
								60,
								"profile history",
							})(w, r)
						},
						func(w http.ResponseWriter, r *http.Request) {
							ghttp.VerifyHeaderKV("Authorization", "Bearer NEW_VALID_ACCESS_TOKEN")(w, r)
							ghttp.VerifyRequest("GET", "/")(w, r)
							ghttp.RespondWith(200, "RESPONSE")(w, r)
						},
					)
				})

				AfterEach(func() {
					server.Close()
				})

				It("Should try to refresh token and then perform the GET request", func() {
					reader, err := oauth2Config.Get("/", nil)
					Expect(server.ReceivedRequests()).To(HaveLen(2))
					Expect(err).ToNot(HaveOccurred())

					output, _ := ioutil.ReadAll(reader)
					Expect(string(output)).To(Equal("RESPONSE"))
					Expect(oauth2Config.AccessToken.AccessToken).To(Equal("NEW_VALID_ACCESS_TOKEN"))
				})
			})

			Context("access token has expired and no refresh token is present", func() {
				var server *ghttp.Server

				BeforeEach(func() {
					server = ghttp.NewServer()

					oauth2Config.AccessToken = &oauth2.Token{Expiry: time.Now().Add(-time.Minute), AccessToken: "InValid-Access-Token"}
					oauth2Config.AuthToken = &AuthToken{AuthCode: "AuthCode"}
					oauth2Config.Endpoint = oauth2.Endpoint{
						TokenURL: server.URL(),
					}
					oauth2Config.RootUrl = server.URL()
				})

				AfterEach(func() {
					server.Close()
				})

				It("Should try to refresh token and then perform the GET request", func() {
					_, err := oauth2Config.Get("/", nil)
					Expect(server.ReceivedRequests()).To(HaveLen(0))
					Expect(err).To(HaveOccurred())
				})
			})
		})

	})

})
