package storage

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/tscolari/cf-broker-api/common/repository"
	"gopkg.in/yaml.v2"
)

func NewLocalFile(location string, capacity int) (*LocalFile, error) {
	var err error
	state := State{
		Capacity:  capacity,
		Instances: map[string]repository.Instance{},
	}

	if _, err = os.Stat(location); os.IsNotExist(err) {
		_, err = os.Create(location)
		if err != nil {
			return nil, err
		}
	}

	localFile := &LocalFile{
		location: location,
		state:    state,
	}

	localFile.Reload()

	return localFile, nil
}

type LocalFile struct {
	location string
	state    State
}

type State struct {
	Capacity  int                            `yaml:"capacity"`
	Instances map[string]repository.Instance `yaml:"instances"`
}

func (s *LocalFile) AvailableInstances() int {
	return s.state.Capacity
}

func (s *LocalFile) InstanceExists(instanceID string) bool {
	if _, exists := s.state.Instances[instanceID]; exists {
		return true
	}

	return false
}

func (s *LocalFile) Instance(instanceID string) (*repository.Instance, error) {
	if instance, exists := s.state.Instances[instanceID]; exists {
		return &instance, nil
	}

	return nil, errors.New("Instance not found")
}

func (s *LocalFile) AddInstance(instance repository.Instance) error {
	if s.state.Capacity == 0 {
		return errors.New("Can't allocate instance, no capacity")
	}

	if _, exists := s.state.Instances[instance.ID]; exists {
		return errors.New("Instance ID is taken")
	}

	s.state.Capacity--
	s.state.Instances[instance.ID] = instance
	s.Save()
	return nil
}

func (s *LocalFile) UpdateInstance(instance repository.Instance) error {
	if _, exists := s.state.Instances[instance.ID]; !exists {
		return errors.New("Instance not found")
	}

	s.state.Instances[instance.ID] = instance
	s.Save()
	return nil
}

func (s *LocalFile) DeleteInstance(instanceID string) error {
	if _, exists := s.state.Instances[instanceID]; !exists {
		return errors.New("Instance not found")
	}

	s.state.Capacity++
	delete(s.state.Instances, instanceID)
	s.Save()
	return nil
}

func (s *LocalFile) InstanceBindingExists(instanceID, bindingID string) bool {
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

func (s *LocalFile) AddInstanceBinding(instanceID, bindingID string) error {
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

func (s *LocalFile) DeleteInstanceBinding(instanceID, bindingID string) error {
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

func (s *LocalFile) Save() error {
	rawData, err := yaml.Marshal(s.state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.location, rawData, 0600)
}

func (s *LocalFile) Reload() error {
	rawData, err := ioutil.ReadFile(s.location)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(rawData, &s.state)
}
