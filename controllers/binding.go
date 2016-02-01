package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/cf-broker-api/common/repository"
	"github.com/tscolari/memcached-broker/app"
)

type Binding struct {
	goa.Controller
	state repository.State
}

func NewBinding(state repository.State) *Binding {
	return &Binding{
		state: state,
	}
}

func (b *Binding) Update(ctx *app.UpdateBindingContext) error {
	if !b.state.InstanceExists(ctx.InstanceId) {
		return ctx.NotFound()
	}

	if b.state.InstanceBindingExists(ctx.InstanceId, ctx.BindingId) {
		return ctx.Conflict()
	}

	err := b.state.AddInstanceBinding(ctx.InstanceId, ctx.BindingId)
	if err != nil {
		return ctx.InternalServerError()
	}

	return ctx.Created()
}

func (b *Binding) Delete(ctx *app.DeleteBindingContext) error {
	if !b.state.InstanceExists(ctx.InstanceId) {
		return ctx.Gone()
	}

	if !b.state.InstanceBindingExists(ctx.InstanceId, ctx.BindingId) {
		return ctx.Gone()
	}

	err := b.state.DeleteInstanceBinding(ctx.InstanceId, ctx.BindingId)
	if err != nil {
		return ctx.InternalServerError()
	}

	return ctx.OK(&app.CfbrokerDashboard{})
}
