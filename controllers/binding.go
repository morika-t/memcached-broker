package controllers

import (
	"github.com/raphael/goa"
	"github.com/tscolari/memcached-broker/app"
	"github.com/tscolari/memcached-broker/config"
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
	var instance config.Instance
	var exists bool
	state := b.storage.GetState()

	if instance, exists = state.Instances[ctx.InstanceId]; !exists {
		return ctx.NotFound()
	}

	for _, binding := range instance.Bindings {
		if binding == ctx.BindingId {
			return ctx.Conflict()
		}
	}

	instance.Bindings = append(instance.Bindings, ctx.BindingId)
	state.Instances[ctx.InstanceId] = instance

	b.storage.PutState(state)
	b.storage.Save()

	return ctx.Created()
}

func (b *Binding) Delete(ctx *app.DeleteBindingContext) error {
	var exists bool
	var instance config.Instance
	state := b.storage.GetState()

	if instance, exists = state.Instances[ctx.InstanceId]; !exists {
		return ctx.Gone()
	}

	for i, binding := range instance.Bindings {
		if binding == ctx.BindingId {
			instance.Bindings = append(instance.Bindings[:i], instance.Bindings[i+1:]...)
			state.Instances[ctx.InstanceId] = instance

			b.storage.PutState(state)
			b.storage.Save()

			return ctx.OK(&app.CfbrokerDashboard{})
		}
	}

	return ctx.Gone()
}
