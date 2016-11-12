package api_test

import (
	. "github.com/DennisDenuto/uber-sdk/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/DennisDenuto/uber-sdk/client/clientfakes"
	. "github.com/tjarratt/gcounterfeiter"
	"fmt"
	"strings"
)

var _ = Describe("Products", func() {

	Describe("GetProduct", func() {
		Context("User has provided valid server token", func() {
			var client ProductsInfo
			var fakeServerTokenClient *clientfakes.FakeClient

			Context("Uber service returns valid successful response", func() {
				BeforeEach(func() {
					fakeServerTokenClient = &clientfakes.FakeClient{}

					fakeServerTokenClient.GetReturns(strings.NewReader(`{
										   "capacity": 4,
										   "description": "The original Uber",
										   "price_details": {
										      "distance_unit": "mile",
										      "cost_per_minute": 0.65,
										      "service_fees": [],
										      "minimum": 15.0,
										      "cost_per_distance": 3.75,
										      "base": 8.0,
										      "cancellation_fee": 10.0,
										      "currency_code": "USD"
										   },
										   "image": "http: //d1a3f4spazzrp4.cloudfront.net/car.jpg",
										   "display_name": "UberBLACK",
										   "product_id": "d4abaae7-f4d6-4152-91cc-77523e8165a4",
										   "shared": false
										}`), nil)
					client = Products{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should generate valid time estimate returned by uber", func() {
					resp, err := client.GetProduct("123")

					Expect(err).ToNot(HaveOccurred())
					Expect(resp).ToNot(BeNil())

					Expect(fakeServerTokenClient).
					To(HaveReceived("Get").
					With(Equal("v1/products/123")).AndWith(BeNil()))

					Expect(resp.PriceDetails.DistanceUnit).To(Equal("mile"))
					Expect(resp.Description).To(Equal("The original Uber"))
					Expect(resp.Capacity).To(Equal(4))
					Expect(resp.Image).To(Equal("http: //d1a3f4spazzrp4.cloudfront.net/car.jpg"))
					Expect(resp.Shared).To(Equal(false))
				})

				Context("Invalid json format", func() {
					BeforeEach(func() {
						fakeServerTokenClient = &clientfakes.FakeClient{}
						fakeServerTokenClient.GetReturns(strings.NewReader("NOT VALID JSON"), nil)

						client = Products{
							ServerTokenClient: fakeServerTokenClient,
						}
					})

					It("Should return an error", func() {
						_, err := client.GetProduct("321")

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Unable to parse Product response from uber: invalid character 'N' looking for beginning of value"))
					})
				})
			})

			Context("Uber service returns non-successful response", func() {
				var uberError error

				BeforeEach(func() {
					uberError = fmt.Errorf("Some random error occurred!")

					fakeServerTokenClient = &clientfakes.FakeClient{}
					fakeServerTokenClient.GetReturns(nil, uberError)

					client = Products{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should return error message returned by uber", func() {
					_, err := client.GetProduct("321")

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("Unable to get Product Info from uber: Some random error occurred!"))
				})
			})

		})
	})

	Describe("ListProducts", func() {
		Context("User has provided valid server token", func() {
			var client ProductsInfo
			var fakeServerTokenClient *clientfakes.FakeClient

			Context("Uber service returns valid successful response", func() {
				BeforeEach(func() {
					fakeServerTokenClient = &clientfakes.FakeClient{}

					fakeServerTokenClient.GetReturns(strings.NewReader(`{
											  "products": [
											    {
											       "capacity": 2,
											       "description": "Ride for less with uberPOOL",
											       "short_description": "POOL",
											       "price_details": null,
											       "cash_enabled": false,
											       "image": "http://d1a3f4spazzrp4.cloudfront.net/car.jpg",
											       "display_name": "POOL",
											       "product_id": "26546650-e557-4a7b-86e7-6a3942445247",
											       "shared": true
											    },
											    {
											       "capacity": 4,
											       "description": "The low-cost Uber",
											       "short_description": "uberX",
											       "price_details": {
												  "distance_unit": "mile",
												  "cost_per_minute": 0.26,
												  "service_fees": [
												     {
													"fee": 1.0,
													"name": "Safe Rides Fee"
												     }
												  ],
												  "minimum": 5.0,
												  "cost_per_distance": 1.3,
												  "base": 2.2,
												  "cancellation_fee": 5.0,
												  "currency_code": "USD"
											       },
											       "cash_enabled": false,
											       "image": "http://d1a3f4spazzrp4.cloudfront.net/car.jpg",
											       "display_name": "uberX",
											       "product_id": "a1111c8c-c720-46c3-8534-2fcdd730040d",
											       "shared": false
											    }
											  ]
											}`), nil)
					client = Products{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should generate valid time estimate returned by uber", func() {
					resp, err := client.ListProducts("123", "321")

					Expect(err).ToNot(HaveOccurred())
					Expect(resp).ToNot(BeNil())

					Expect(fakeServerTokenClient).
					To(HaveReceived("Get").
					With(Equal("v1/products")).AndWith(BeNil()))

					Expect(resp.Products).To(HaveLen(2))
				})

				Context("Invalid json format", func() {
					BeforeEach(func() {
						fakeServerTokenClient = &clientfakes.FakeClient{}
						fakeServerTokenClient.GetReturns(strings.NewReader("NOT VALID JSON"), nil)

						client = Products{
							ServerTokenClient: fakeServerTokenClient,
						}
					})

					It("Should return an error", func() {
						_, err := client.ListProducts("321", "123")

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Unable to parse Products response from uber: invalid character 'N' looking for beginning of value"))
					})
				})
			})

			Context("Uber service returns non-successful response", func() {
				var uberError error

				BeforeEach(func() {
					uberError = fmt.Errorf("Some random error occurred!")

					fakeServerTokenClient = &clientfakes.FakeClient{}
					fakeServerTokenClient.GetReturns(nil, uberError)

					client = Products{
						ServerTokenClient: fakeServerTokenClient,
					}
				})

				It("should return error message returned by uber", func() {
					_, err := client.ListProducts("321", "123")

					Expect(err).To(HaveOccurred())

					fmt.Println(err.Error())
					Expect(err).To(MatchError("Unable to get Products from uber: Some random error occurred!"))
				})
			})

		})
	})
})
