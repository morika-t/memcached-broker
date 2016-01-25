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

var _ = Describe("Binding", func() {
	var bindingController *controllers.Binding
	var storage *fakes.FakeStorage
	var goaContext *goa.Context
	var responseWriter *httptest.ResponseRecorder

	BeforeEach(func() {
		storage = new(fakes.FakeStorage)
		bindingController = controllers.NewBinding(storage)

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
			bindingContext.InstanceId = "instance-1"
			bindingContext.BindingId = "binding-1"
			bindingContext.AppGuid = "app-guid"

			Expect(err).ToNot(HaveOccurred())
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				state := config.State{
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							ID: "instance-1",
						},
					},
				}

				storage.GetStateReturns(state)
				err := bindingController.Update(bindingContext)
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
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							ID:       "instance-1",
							Bindings: []string{"binding-1"},
						},
					},
				}))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state := config.State{
					Instances: map[string]config.Instance{},
				}

				storage.GetStateReturns(state)
				err := bindingController.Update(bindingContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 404", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(404))
			})
		})

		Context("when the binding id already exists", func() {
			BeforeEach(func() {
				state := config.State{
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							Bindings: []string{"binding-1"},
						},
					},
				}

				storage.GetStateReturns(state)
				err := bindingController.Update(bindingContext)
				Expect(err).ToNot(HaveOccurred())
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
			bindingContext.InstanceId = "instance-1"
			bindingContext.BindingId = "binding-1"

			Expect(err).ToNot(HaveOccurred())
		})

		Context("when all goes ok", func() {
			BeforeEach(func() {
				state := config.State{
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							ID:       "instance-1",
							Bindings: []string{"binding-1"},
						},
					},
				}

				storage.GetStateReturns(state)
				err := bindingController.Delete(bindingContext)
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
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							ID:       "instance-1",
							Bindings: []string{},
						},
					},
				}))
			})
		})

		Context("when the instance doesn't exist", func() {
			BeforeEach(func() {
				state := config.State{
					Instances: map[string]config.Instance{},
				}

				storage.GetStateReturns(state)
				err := bindingController.Delete(bindingContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 410", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(410))
			})
		})

		Context("when the binding doesn't exist", func() {
			BeforeEach(func() {
				state := config.State{
					Capacity: 1,
					Instances: map[string]config.Instance{
						"instance-1": config.Instance{
							ID:       "instance-1",
							Bindings: []string{},
						},
					},
				}

				storage.GetStateReturns(state)
				err := bindingController.Delete(bindingContext)
				Expect(err).ToNot(HaveOccurred())
			})

			It("responds with 410", func() {
				Expect(goaContext.ResponseStatus()).To(Equal(410))
			})
		})
	})
})
