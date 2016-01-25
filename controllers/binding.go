package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/storage"
)

type Binding struct {
	goa.Controller
	storage storage.Storage
}

func NewBinding(storage storage.Storage) *Binding {
	return &Binding{
		storage: storage,
	}
}

func (b *Binding) Update(ctx *app.UpdateBindingContext) error {
	state := b.storage.GetState()

	if !state.InstanceExists(ctx.InstanceId) {
		return ctx.NotFound()
	}

	if state.InstanceBindingExists(ctx.InstanceId, ctx.BindingId) {
		return ctx.Conflict()
	}

	err := state.AddInstanceBinding(ctx.InstanceId, ctx.BindingId)
	if err != nil {
		return ctx.InternalServerError()
	}

	b.storage.PutState(state)
	b.storage.Save()
	return ctx.Created()
}

func (b *Binding) Delete(ctx *app.DeleteBindingContext) error {
	state := b.storage.GetState()

	if !state.InstanceExists(ctx.InstanceId) {
		return ctx.Gone()
	}

	if !state.InstanceBindingExists(ctx.InstanceId, ctx.BindingId) {
		return ctx.Gone()
	}

	err := state.DeleteInstanceBinding(ctx.InstanceId, ctx.BindingId)
	if err != nil {
		return ctx.InternalServerError()
	}

	b.storage.PutState(state)
	b.storage.Save()

	return ctx.OK(&app.CfbrokerDashboard{})
}
