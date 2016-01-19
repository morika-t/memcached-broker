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

	if state.InstanceExists(ctx.InstanceId) {
		return ctx.Conflict()
	}

	instance := config.Instance{
		ServiceID:      ctx.ServiceId,
		PlanID:         ctx.PlanId,
		OrganizationID: ctx.OrganizationId,
		SpaceID:        ctx.SpaceId,
	}

	err := state.AddInstance(ctx.InstanceId, instance)
	if err != nil {
		return ctx.ServiceUnavailable()
	}

	p.storage.PutState(state)
	p.storage.Save()

	return ctx.Created()
}

func (p *Provisioning) Update(ctx *app.UpdateProvisioningContext) error {
	state := p.storage.GetState()

	instance, err := state.Instance(ctx.InstanceId)
	if err != nil {
		return ctx.NotFound()
	}

	instance.ServiceID = ctx.ServiceId
	instance.PlanID = ctx.PlanId

	state.UpdateInstance(ctx.InstanceId, *instance)
	p.storage.PutState(state)
	p.storage.Save()

	return ctx.OK(&app.CfbrokerDashboard{})
}

func (p *Provisioning) Delete(ctx *app.DeleteProvisioningContext) error {
	state := p.storage.GetState()

	if !state.InstanceExists(ctx.InstanceId) {
		return ctx.Gone()
	}

	err := state.DeleteInstance(ctx.InstanceId)
	if err != nil {
		return ctx.Gone()
	}

	p.storage.PutState(state)
	p.storage.Save()

	return ctx.OK(&app.CfbrokerDashboard{})
}
