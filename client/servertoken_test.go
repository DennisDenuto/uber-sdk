package client_test

import (
	. "github.com/DennisDenuto/uber-client/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Servertoken", func() {

	var server *ghttp.Server
	var serverTokenClient ServerTokenClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		serverTokenClient = ServerTokenClient{
			RootUrl: server.URL(),
		}
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})

	Context("A valid server token is provided", func() {

		Context("Accessing an endpoint with no url params", func() {
			BeforeEach(func() {
				serverTokenClient.ServerToken = "PROVIDED_SERVER_TOKEN"
				server.AppendHandlers(
					ghttp.VerifyRequest("GET", "/foo", "server_token=PROVIDED_SERVER_TOKEN"),
				)
			})

			It("Should set the server token as a url param correctly", func() {
				_, err := serverTokenClient.Get("/foo", nil)

				Expect(err).ToNot(HaveOccurred())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("Accessing an endpoint with query params", func() {
			BeforeEach(func() {
				serverTokenClient.ServerToken = "PROVIDED_SERVER_TOKEN"
				server.AppendHandlers(
					ghttp.VerifyRequest("GET", "/bar", "custom_param=foo&server_token=PROVIDED_SERVER_TOKEN"),
				)
			})

			It("Should set the server token as a url param correctly", func() {
				var customParams map[string]string = make(map[string]string, 1)
				customParams["custom_param"] = "foo"

				_, err := serverTokenClient.Get("/bar", customParams)

				Expect(err).ToNot(HaveOccurred())
				Expect(server.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})
})
