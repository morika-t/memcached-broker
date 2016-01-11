package main

import (
	"github.com/raphael/goa"
	"github.com/raphael/goa/examples/cellar/swagger"
	"github.com/tscolari/cf-broker-api/common/objects"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/controllers"
)

func main() {
	service := goa.New("cfbroker")

	catalogController := controllers.NewCatalog(objects.Catalog{})
	app.MountCatalogController(service, catalogController)

	swagger.MountController(service)
	service.ListenAndServe(":8080")
}
