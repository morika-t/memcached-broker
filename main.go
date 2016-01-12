package main

import (
	"github.com/raphael/goa"
	"github.com/raphael/goa/examples/cellar/swagger"
	"github.com/tscolari/cf-broker-api/common/controllers"
	"github.com/tscolari/memcached-broker/app"
)

func main() {
	service := goa.New("cfbroker")

	catalogController := controllers.NewCatalog(app.CfbrokerCatalog{})
	app.MountCatalogController(service, catalogController)

	swagger.MountController(service)
	service.ListenAndServe(":8080")
}
