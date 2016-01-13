package storage

import "github.com/tscolari/memcached-broker/config"

type Storage interface {
	Save() error
	Reload() error

	GetState() config.State
	PutState(state config.State)
}
