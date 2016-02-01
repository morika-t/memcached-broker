package controllers_test

import (
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

var _ = Describe("Binding", func() {
	var bindingController *controllers.Binding
	var state *fakes.FakeState
	var goaContext *goa.Context
	var responseWriter *httptest.ResponseRecorder

	BeforeEach(func() {
		state = new(fakes.FakeState)
		gctx := context.Background()
		req := http.Request{}
		responseWriter = httptest.NewRecorder()
		params := url.Values{}
		payload := map[string]string{}

		goaContext = goa.NewContext(gctx, &req, responseWriter, params, payload)
	})

	Describe("#Update", func() {
		var bindingContext *app.UpdateBindingContext

		BeforeEach(func() {
			var err error
			bindingContext, err = app.NewUpdateBindingContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			bindingContext.InstanceId = "instance-1"
			bindingContext.BindingId = "binding-1"
			bindingContext.AppGuid = "app-guid"
		})

		JustBeforeEach(func() {
			bindingController = controllers.NewBinding(state)
			err := bindingController.Update(bindingContext)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when all goes ok", func() {

			BeforeEach(func() {
				instance := repository.Instance{
					ID: "instance-1",
				}

				state.InstanceExistsReturns(true)
				state.InstanceBindingExistsReturns(false)
				state.InstanceReturns(&instance, nil)
			})

			It("responds with 201", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(201))
			})

			It("updates the state", func() {
				instanceID, bindingID := state.AddInstanceBindingArgsForCall(0)
				Expect(instanceID).To(Equal("instance-1"))
				Expect(bindingID).To(Equal("binding-1"))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
			})

			It("responds with 404", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(404))
			})
		})

		Context("when the binding id already exists", func() {
			BeforeEach(func() {
				state.InstanceExistsReturns(true)
				state.InstanceBindingExistsReturns(true)
			})

			It("responds with 409", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(409))
			})
		})
	})

	Describe("#Delete", func() {
		var bindingContext *app.DeleteBindingContext

		BeforeEach(func() {
			var err error
			bindingContext, err = app.NewDeleteBindingContext(goaContext)
			Expect(err).ToNot(HaveOccurred())

			bindingContext.InstanceId = "instance-1"
			bindingContext.BindingId = "binding-1"
		})

		JustBeforeEach(func() {
			bindingController = controllers.NewBinding(state)
			err := bindingController.Delete(bindingContext)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				instance := repository.Instance{
					ID:       "instance-1",
					Bindings: []string{"binding-1"},
				}

				state.InstanceExistsReturns(true)
				state.InstanceBindingExistsReturns(true)
				state.InstanceReturns(&instance, nil)
			})

			It("responds with 200", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(200))
			})

			It("sends the delete message to the state", func() {
				instanceID, bindingID := state.DeleteInstanceBindingArgsForCall(0)
				Expect(instanceID).To(Equal("instance-1"))
				Expect(bindingID).To(Equal("binding-1"))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state.InstanceExistsReturns(false)
			})

			It("responds with 410", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(410))
			})
		})

		Context("when the binding doesn't exist", func() {
			BeforeEach(func() {
				state.InstanceExistsReturns(true)
				state.InstanceBindingExistsReturns(false)
			})

			It("responds with 410", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(410))
			})
		})
	})
})
