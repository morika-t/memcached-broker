package config

import "errors"

func NewState(capacity int) State {
	return State{
		Capacity:  capacity,
		Instances: map[string]Instance{},
	}
}

type State struct {
	Capacity  int                 `yaml:"capacity"`
	Instances map[string]Instance `yaml:"instances"`
}

func (s *State) InstanceExists(instanceID string) bool {
	if _, exists := s.Instances[instanceID]; exists {
		return true
	}

	return false
}

func (s *State) Instance(instanceID string) (*Instance, error) {
	if instance, exists := s.Instances[instanceID]; exists {
		return &instance, nil
	}

	return nil, errors.New("Instance not found")
}

func (s *State) AddInstance(instanceID string, instance Instance) error {
	if s.Capacity == 0 {
		return errors.New("Can't allocate instance, no capacity")
	}

	if _, exists := s.Instances[instanceID]; exists {
		return errors.New("Instance ID is taken")
	}

	s.Capacity--
	s.Instances[instanceID] = instance
	return nil
}

func (s *State) UpdateInstance(instanceID string, instance Instance) error {
	if _, exists := s.Instances[instanceID]; !exists {
		return errors.New("Instance not found")
	}

	s.Instances[instanceID] = instance
	return nil
}

func (s *State) DeleteInstance(instanceID string) error {
	if _, exists := s.Instances[instanceID]; !exists {
		return errors.New("Instance not found")
	}

	s.Capacity++
	delete(s.Instances, instanceID)
	return nil
}
