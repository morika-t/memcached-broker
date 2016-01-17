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
	var goaContext *goa.Context
	var responseWriter *httptest.ResponseRecorder

	BeforeEach(func() {
		storage = new(fakes.FakeStorage)
		provisioningController = controllers.NewProvisioning(storage)

		gctx := context.Background()
		req := http.Request{}
		responseWriter = httptest.NewRecorder()
		params := url.Values{}
		payload := map[string]string{}

		goaContext = goa.NewContext(gctx, &req, responseWriter, params, payload)
	})

	Describe("#Create", func() {
		var provisioningContext *app.CreateProvisioningContext

		BeforeEach(func() {
			var err error
			provisioningContext, err = app.NewCreateProvisioningContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			provisioningContext.InstanceId = "some-instance-id"
			provisioningContext.OrganizationId = "org-1"
			provisioningContext.SpaceId = "space-1"
			provisioningContext.ServiceId = "service-1"
			provisioningContext.PlanId = "plan-1"
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
						"some-instance-id": config.Instance{
							ServiceID:      "service-1",
							PlanID:         "plan-1",
							OrganizationID: "org-1",
							SpaceID:        "space-1",
						},
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

	Describe("#Update", func() {
		var provisioningContext *app.UpdateProvisioningContext

		BeforeEach(func() {
			var err error
			provisioningContext, err = app.NewUpdateProvisioningContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			provisioningContext.InstanceId = "some-instance-id"
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity: 0,
					Instances: map[string]config.Instance{
						"some-instance-id": config.Instance{
							ServiceID:      "service-1",
							PlanID:         "plan-1",
							OrganizationID: "org-1",
							SpaceID:        "space-1",
						},
					},
				}

				storage.GetStateReturns(state)
				provisioningContext.ServiceId = "service-2"
				provisioningContext.PlanId = "plan-2"
				err := provisioningController.Update(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 200", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(200))
			})

			It("updates the state file with correct information", func() {
				Expect(storage.PutStateCallCount()).To(Equal(1))
				Expect(storage.SaveCallCount()).To(Equal(1))

				receivedState := storage.PutStateArgsForCall(0)
				Expect(receivedState).To(Equal(config.State{
					Capacity: 0,
					Instances: map[string]config.Instance{
						"some-instance-id": config.Instance{
							ServiceID:      "service-2",
							PlanID:         "plan-2",
							OrganizationID: "org-1",
							SpaceID:        "space-1",
						},
					},
				}))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity:  1,
					Instances: map[string]config.Instance{},
				}

				storage.GetStateReturns(state)
				err := provisioningController.Update(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 404", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(404))
			})
		})
	})

	Describe("#Delete", func() {
		var provisioningContext *app.DeleteProvisioningContext

		BeforeEach(func() {
			var err error
			provisioningContext, err = app.NewDeleteProvisioningContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			provisioningContext.InstanceId = "some-instance-id"
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity: 0,
					Instances: map[string]config.Instance{
						"some-instance-id": config.Instance{},
					},
				}

				storage.GetStateReturns(state)
				err := provisioningController.Delete(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 200", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(200))
			})

			It("deletes the instance from the state file", func() {
				Expect(storage.PutStateCallCount()).To(Equal(1))
				Expect(storage.SaveCallCount()).To(Equal(1))

				receivedState := storage.PutStateArgsForCall(0)
				Expect(receivedState).To(Equal(config.State{
					Capacity:  1,
					Instances: map[string]config.Instance{},
				}))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity:  1,
					Instances: map[string]config.Instance{},
				}

				storage.GetStateReturns(state)
				err := provisioningController.Delete(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 404", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(404))
			})
		})
	})
})
