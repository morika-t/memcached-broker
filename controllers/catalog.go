package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/cf-broker-api/common/objects"
	"github.com/tscolari/memcached-broker/app"
)

type Catalog struct {
	goa.Controller
	catalog objects.Catalog
}

type ShowCatalogContext struct {
	*goa.Context
}

func NewShowCatalogContext(ctx *goa.Context) (*ShowCatalogContext, error) {
	return &ShowCatalogContext{ctx}, nil
}

func NewCatalog(config objects.Catalog) *Catalog {
	return &Catalog{
		catalog: config,
	}
}

func (c *Catalog) Show(ctx *app.ShowCatalogContext) error {
	return nil
}
