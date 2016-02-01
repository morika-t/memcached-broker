package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/cf-broker-api/common/repository"
	"github.com/tscolari/memcached-broker/app"
)

type Provisioning struct {
	goa.Controller
	state repository.State
}

func NewProvisioning(state repository.State) *Provisioning {
	return &Provisioning{
		state: state,
	}
}

func (p *Provisioning) Create(ctx *app.CreateProvisioningContext) error {
	if p.state.InstanceExists(ctx.InstanceId) {
		return ctx.Conflict()
	}

	instance := repository.Instance{
		ID:             ctx.InstanceId,
		ServiceID:      ctx.ServiceId,
		PlanID:         ctx.PlanId,
		OrganizationID: ctx.OrganizationId,
		SpaceID:        ctx.SpaceId,
	}

	err := p.state.AddInstance(instance)
	if err != nil {
		return ctx.ServiceUnavailable()
	}

	return ctx.Created()
}

func (p *Provisioning) Update(ctx *app.UpdateProvisioningContext) error {
	instance, err := p.state.Instance(ctx.InstanceId)
	if err != nil {
		return ctx.NotFound()
	}

	instance.ServiceID = ctx.ServiceId
	instance.PlanID = ctx.PlanId

	p.state.UpdateInstance(*instance)

	return ctx.OK(&app.CfbrokerDashboard{})
}

func (p *Provisioning) Delete(ctx *app.DeleteProvisioningContext) error {
	if !p.state.InstanceExists(ctx.InstanceId) {
		return ctx.Gone()
	}

	err := p.state.DeleteInstance(ctx.InstanceId)
	if err != nil {
		return ctx.Gone()
	}

	return ctx.OK(&app.CfbrokerDashboard{})
}
