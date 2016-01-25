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

func (s *State) AddInstance(instance Instance) error {
	if s.Capacity == 0 {
		return errors.New("Can't allocate instance, no capacity")
	}

	if _, exists := s.Instances[instance.ID]; exists {
		return errors.New("Instance ID is taken")
	}

	s.Capacity--
	s.Instances[instance.ID] = instance
	return nil
}

func (s *State) UpdateInstance(instance Instance) error {
	if _, exists := s.Instances[instance.ID]; !exists {
		return errors.New("Instance not found")
	}

	s.Instances[instance.ID] = instance
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

func (s *State) InstanceBindingExists(instanceID, bindingID string) bool {
	instance, err := s.Instance(instanceID)
	if err != nil {
		return false
	}

	for _, binding := range instance.Bindings {
		if binding == bindingID {
			return true
		}
	}

	return false
}

func (s *State) AddInstanceBinding(instanceID, bindingID string) error {
	instance, err := s.Instance(instanceID)
	if err != nil {
		return err
	}

	for _, binding := range instance.Bindings {
		if binding == bindingID {
			return errors.New("Binding ID is taken")
		}
	}

	instance.Bindings = append(instance.Bindings, bindingID)
	return s.UpdateInstance(*instance)
}

func (s *State) DeleteInstanceBinding(instanceID, bindingID string) error {
	instance, err := s.Instance(instanceID)
	if err != nil {
		return err
	}

	for i, binding := range instance.Bindings {
		if binding == bindingID {
			instance.Bindings = append(instance.Bindings[:i], instance.Bindings[i+1:]...)
			return s.UpdateInstance(*instance)
		}
	}

	return errors.New("Binding not found")
}
