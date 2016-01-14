package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/config"
	"github.com/tscolari/memcached-broker/controllers"
	"github.com/tscolari/memcached-broker/storage/fakes"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioning", func() {
	var provisioningController *controllers.Provisioning
	var storage *fakes.FakeStorage

	BeforeEach(func() {
		storage = new(fakes.FakeStorage)
		provisioningController = controllers.NewProvisioning(storage)
	})

	Describe("#Create", func() {
		var goaContext *goa.Context
		var provisioningContext *app.CreateProvisioningContext
		var responseWriter *httptest.ResponseRecorder

		BeforeEach(func() {
			gctx := context.Background()
			req := http.Request{}
			responseWriter = httptest.NewRecorder()
			params := url.Values{}
			payload := map[string]string{}

			goaContext = goa.NewContext(gctx, &req, responseWriter, params, payload)

			var err error
			provisioningContext, err = app.NewCreateProvisioningContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			provisioningContext.InstanceId = "some-instance-id"
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity:  1,
					Instances: map[string]config.Instance{},
				}

				storage.GetStateReturns(state)
				err := provisioningController.Create(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 201", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(201))
			})

			It("updates the state file", func() {
				Expect(storage.PutStateCallCount()).To(Equal(1))
				Expect(storage.SaveCallCount()).To(Equal(1))

				receivedState := storage.PutStateArgsForCall(0)
				Expect(receivedState).To(Equal(config.State{
					Capacity: 0,
					Instances: map[string]config.Instance{
						"some-instance-id": config.Instance{},
					},
				}))
			})
		})

		Context("when the instance id already exists", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity: 1,
					Instances: map[string]config.Instance{
						"some-instance-id": config.Instance{},
					},
				}

				storage.GetStateReturns(state)
				err := provisioningController.Create(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 409", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(409))
			})
		})

		Context("when there's no capacity", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity: 0,
				}

				storage.GetStateReturns(state)
				err := provisioningController.Create(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 503", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(503))
			})
		})
	})
})
