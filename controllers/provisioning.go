package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/storage"
)

type Provisioning struct {
	goa.Controller
	storage *storage.Storage
}

func NewProvisioning(storage *storage.Storage) *Provisioning {
	return &Provisioning{
		storage: storage,
	}
}

func (p *Provisioning) Create(ctx *app.CreateProvisioningContext) error {
	return nil
}
