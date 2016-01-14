package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/controllers"
	"github.com/tscolari/memcached-broker/storage"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioning", func() {
	var provisioningController *controllers.Provisioning
	var storage storage.Storage

	BeforeEach(func() {
		provisioningController = controllers.NewProvisioning(&storage)
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

			err = provisioningController.Create(provisioningContext)
			Expect(err).ToNot(HaveOccurred())
		})

		It("responds with 201", func() {
			Expect(goaContext.ResponseStatus()).To(Equal(201))
		})

		It("updates the state file", func() {
		})
	})
})
