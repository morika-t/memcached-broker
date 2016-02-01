package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/raphael/goa"
	"github.com/tscolari/cf-broker-api/common/repository"
	"github.com/tscolari/cf-broker-api/common/repository/fakes"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/controllers"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioning", func() {
	var provisioningController *controllers.Provisioning
	var goaContext *goa.Context
	var responseWriter *httptest.ResponseRecorder
	var state *fakes.FakeState

	BeforeEach(func() {
		state = new(fakes.FakeState)
		provisioningController = controllers.NewProvisioning(state)

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
				err := provisioningController.Create(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 201", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(201))
			})

			It("sends the correct message to the state", func() {
				recordedInstance := state.AddInstanceArgsForCall(0)
				Expect(recordedInstance.ID).To(Equal("some-instance-id"))
				Expect(recordedInstance.OrganizationID).To(Equal("org-1"))
				Expect(recordedInstance.SpaceID).To(Equal("space-1"))
				Expect(recordedInstance.ServiceID).To(Equal("service-1"))
				Expect(recordedInstance.PlanID).To(Equal("plan-1"))
			})
		})

		Context("when the instance id already exists", func() {
			BeforeEach(func() {
				state.InstanceExistsReturns(true)

				err := provisioningController.Create(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 409", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(409))
			})
		})

		Context("when there's no capacity", func() {
			BeforeEach(func() {
				state.AddInstanceReturns(errors.New("Failed"))
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
				instance := repository.Instance{
					ID:             "some-instance-id",
					ServiceID:      "service-1",
					PlanID:         "plan-1",
					OrganizationID: "org-1",
					SpaceID:        "space-1",
				}

				state.InstanceReturns(&instance, nil)

				provisioningContext.ServiceId = "service-2"
				provisioningContext.PlanId = "plan-2"
				err := provisioningController.Update(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 200", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(200))
			})

			It("sends the correct message to the state", func() {
				recordedInstance := state.UpdateInstanceArgsForCall(0)
				Expect(recordedInstance.ID).To(Equal("some-instance-id"))
				Expect(recordedInstance.OrganizationID).To(Equal("org-1"))
				Expect(recordedInstance.SpaceID).To(Equal("space-1"))
				Expect(recordedInstance.ServiceID).To(Equal("service-2"))
				Expect(recordedInstance.PlanID).To(Equal("plan-2"))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state.InstanceReturns(nil, errors.New("Not here!"))
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
				state.InstanceExistsReturns(true)

				err := provisioningController.Delete(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 200", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(200))
			})

			It("sends the correct message to the state", func() {
				instanceID := state.DeleteInstanceArgsForCall(0)
				Expect(instanceID).To(Equal("some-instance-id"))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				err := provisioningController.Delete(provisioningContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 410", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(410))
			})
		})
	})
})
