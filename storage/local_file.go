package storage

import (
	"io/ioutil"
	"os"

	"github.com/tscolari/memcached-broker/config"
	"gopkg.in/yaml.v2"
)

func NewLocalFile(location string) (*LocalFile, error) {
	var err error
	state := config.State{}

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
	state    config.State
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

func (s *LocalFile) GetState() config.State {
	return s.state
}

func (s *LocalFile) PutState(state config.State) {
	s.state = state
}
