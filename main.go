package main

import (
	"github.com/raphael/goa"
	"github.com/raphael/goa/examples/cellar/swagger"
	commoncontrollers "github.com/tscolari/cf-broker-api/common/controllers"
	"github.com/tscolari/cf-broker-api/common/storage"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/config"
	"github.com/tscolari/memcached-broker/controllers"
)

func main() {
	service := goa.New("cfbroker")
	config, err := config.Load("./config.yaml")
	if err != nil {
		panic(err)
	}

	stateFile, err := storage.NewYamlFile(config.StateFile)
	if err != nil {
		panic(err)
	}

	provisioningController := controllers.NewProvisioning(stateFile)

	catalogController := commoncontrollers.NewCatalog(config.Catalog)
	app.MountCatalogController(service, catalogController)
	app.MountProvisioningController(service, provisioningController)

	swagger.MountController(service)
	service.ListenAndServe(":8080")
}
