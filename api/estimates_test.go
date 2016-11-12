package api_test

import (
	. "github.com/DennisDenuto/uber-client/api"

	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"strings"
)

type FakeServerTokenClient struct {
	ExpectedUrl         string
	ExpectedQueryParams map[string]string
	ReturnedResponse    string
	ReturnedError       error
}

func (fakeClient *FakeServerTokenClient) Get(url string, queryParams map[string]string) (io.Reader, error) {
	fakeClient.ExpectedUrl = url
	fakeClient.ExpectedQueryParams = queryParams

	if fakeClient.ReturnedError != nil {
		return nil, fakeClient.ReturnedError
	}

	timeJsonResponse := strings.NewReader(fakeClient.ReturnedResponse)
	return timeJsonResponse, nil
}

var _ = Describe("Estimates", func() {

	Describe("Price", func() {
		Context("User has provided valid server token", func() {
			var client Estimator
			var fakeServerTokenClient *FakeServerTokenClient

			Context("Uber service returns valid successful response", func() {
				BeforeEach(func() {
					fakeServerTokenClient = &FakeServerTokenClient{
						ReturnedResponse: `
							{
							  "prices": [
							    {
							      "product_id": "26546650-e557-4a7b-86e7-6a3942445247",
							      "currency_code": "USD",
							      "localized_display_name": "POOL",
							      "display_name": "POOL",
							      "estimate": "$5.75",
							      "minimum": null,
							      "low_estimate": 5,
							      "high_estimate": 6,
							      "surge_multiplier": 1,
							      "duration": 640,
							      "distance": 5.34
							    },
							    {
							      "product_id": "08f17084-23fd-4103-aa3e-9b660223934b",
							      "currency_code": "USD",
							      "localized_display_name": "UberBLACK",
							      "display_name": "UberBLACK",
							      "estimate": "$23-29",
							      "minimum": 9,
							      "low_estimate": 23,
							      "high_estimate": 29,
							      "surge_multiplier": 1,
							      "duration": 640,
							      "distance": 5.34
							    }
							  ]
							}`,
					}
					client = Estimate{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should generate valid time estimate returned by uber", func() {
					resp, err := client.GetPrice("1.111111111", "1.22222", "2.222222222", "3.3333", 0)

					Expect(err).ToNot(HaveOccurred())
					Expect(resp).ToNot(BeNil())
					Expect(resp.Prices).To(HaveLen(2))

					Expect(fakeServerTokenClient.ExpectedUrl).To(Equal("v1/estimates/price"))
					Expect(fakeServerTokenClient.ExpectedQueryParams).To(HaveKeyWithValue("start_latitude", "2.222222222"))
					Expect(fakeServerTokenClient.ExpectedQueryParams).To(HaveKeyWithValue("start_longitude", "1.111111111"))
				})

				Context("Invalid json format", func() {
					BeforeEach(func() {
						fakeServerTokenClient = &FakeServerTokenClient{
							ReturnedResponse: `
							NOT VALID JSON`,
						}
						client = Estimate{
							ServerTokenClient: fakeServerTokenClient,
						}
					})

					It("Should return an error", func() {
						_, err := client.GetPrice("1.111111111", "1.22222", "2.222222222", "3.3333", 0)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Unable to parse Price Estimates response from uber: invalid character 'N' looking for beginning of value"))
					})
				})
			})

			Context("Uber service returns non-successful response", func() {
				var uberError error

				BeforeEach(func() {
					uberError = fmt.Errorf("Some random error occurred!")

					fakeServerTokenClient = &FakeServerTokenClient{
						ReturnedError: uberError,
					}

					client = Estimate{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should return error message returned by uber", func() {
					_, err := client.GetPrice("1.111111111", "1.22222", "2.222222222", "3.3333", 0)

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("Unable to get Price Estimates from uber: Some random error occurred!"))
				})
			})

		})
	})

	Describe("Time", func() {
		Context("User has provided valid server token", func() {
			var client Estimator
			var fakeServerTokenClient *FakeServerTokenClient

			Context("Uber service returns valid successful response", func() {
				BeforeEach(func() {
					fakeServerTokenClient = &FakeServerTokenClient{
						ReturnedResponse: `
							{
							   "times": [
							      {
								 "localized_display_name":"uberPOOL",
								 "estimate":180,
								 "display_name":"uberPOOL",
								 "product_id":"26546650-e557-4a7b-86e7-6a3942445247"
							      }
							   ]
							}`,
					}
					client = Estimate{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should generate valid time estimate returned by uber", func() {
					resp, err := client.GetTime("1.111111111", "2.222222222")

					Expect(err).ToNot(HaveOccurred())
					Expect(resp).ToNot(BeNil())
					Expect(resp.Times).To(HaveLen(1))
					Expect(resp.Times[0].DisplayName).To(Equal("uberPOOL"))
					Expect(fakeServerTokenClient.ExpectedUrl).To(Equal("v1/estimates/time"))
					Expect(fakeServerTokenClient.ExpectedQueryParams).To(HaveKeyWithValue("start_latitude", "2.222222222"))
					Expect(fakeServerTokenClient.ExpectedQueryParams).To(HaveKeyWithValue("start_longitude", "1.111111111"))
				})

				Context("Invalid json format", func() {
					BeforeEach(func() {
						fakeServerTokenClient = &FakeServerTokenClient{
							ReturnedResponse: `
							NOT VALID JSON`,
						}
						client = Estimate{
							ServerTokenClient: fakeServerTokenClient,
						}
					})

					It("Should return an error", func() {
						_, err := client.GetTime("1.0", "2.0")

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Unable to parse Time Estimates response from uber: invalid character 'N' looking for beginning of value"))
					})
				})
			})

			Context("Uber service returns non-successful response", func() {
				var uberError error

				BeforeEach(func() {
					uberError = fmt.Errorf("Some random error occurred!")

					fakeServerTokenClient = &FakeServerTokenClient{
						ReturnedError: uberError,
					}

					client = Estimate{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should return error message returned by uber", func() {
					_, err := client.GetTime("1.0", "2.0")

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("Unable to get Time Estimates from uber: Some random error occurred!"))
				})
			})
		})
	})

})
