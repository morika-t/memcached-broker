package main

import (
	"github.com/raphael/goa"
	"github.com/raphael/goa/examples/cellar/swagger"
	commoncontrollers "github.com/tscolari/cf-broker-api/common/controllers"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/config"
	"github.com/tscolari/memcached-broker/controllers"
	"github.com/tscolari/memcached-broker/storage"
)

func main() {
	service := goa.New("cfbroker")
	configuration, err := config.Load("./config.yaml")
	if err != nil {
		panic(err)
	}

	store, err := storage.NewLocalFile("/tmp/data")
	if err != nil {
		panic(err)
	}

	provisioningController := controllers.NewProvisioning(store)
	catalogController := commoncontrollers.NewCatalog(configuration.Catalog)

	app.MountCatalogController(service, catalogController)
	app.MountProvisioningController(service, provisioningController)

	swagger.MountController(service)
	service.ListenAndServe(":8080")
}
