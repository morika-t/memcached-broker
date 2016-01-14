package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/config"
	"github.com/tscolari/memcached-broker/storage"
)

type Provisioning struct {
	goa.Controller
	storage storage.Storage
}

func NewProvisioning(storage storage.Storage) *Provisioning {
	return &Provisioning{
		storage: storage,
	}
}

func (p *Provisioning) Create(ctx *app.CreateProvisioningContext) error {
	state := p.storage.GetState()

	if state.Capacity <= 0 {
		return ctx.ServiceUnavailable()
	}

	if _, exists := state.Instances[ctx.InstanceId]; exists {
		return ctx.Conflict()
	}

	instance := config.Instance{}

	state.Instances[ctx.InstanceId] = instance
	state.Capacity = state.Capacity - 1
	p.storage.PutState(state)
	p.storage.Save()

	return ctx.Created()
}
