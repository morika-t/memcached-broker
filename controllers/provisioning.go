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

	instance := config.Instance{
		ServiceID:      ctx.ServiceId,
		PlanID:         ctx.PlanId,
		OrganizationID: ctx.OrganizationId,
		SpaceID:        ctx.SpaceId,
	}

	state.Instances[ctx.InstanceId] = instance
	state.Capacity = state.Capacity - 1
	p.storage.PutState(state)
	p.storage.Save()

	return ctx.Created()
}

func (p *Provisioning) Update(ctx *app.UpdateProvisioningContext) error {
	var instance config.Instance
	var exists bool
	state := p.storage.GetState()

	if instance, exists = state.Instances[ctx.InstanceId]; !exists {
		return ctx.NotFound()
	}

	instance.ServiceID = ctx.ServiceId
	instance.PlanID = ctx.PlanId

	state.Instances[ctx.InstanceId] = instance
	p.storage.PutState(state)
	p.storage.Save()

	return ctx.OK(&app.CfbrokerDashboard{})
}

func (p *Provisioning) Delete(ctx *app.DeleteProvisioningContext) error {
	var exists bool
	state := p.storage.GetState()

	if _, exists = state.Instances[ctx.InstanceId]; !exists {
		return ctx.Gone()
	}

	state.Capacity++
	delete(state.Instances, ctx.InstanceId)

	p.storage.PutState(state)
	p.storage.Save()

	return ctx.OK(&app.CfbrokerDashboard{})
}
