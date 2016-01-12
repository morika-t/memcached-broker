package config

import (
	"fmt"
	"io/ioutil"

	"github.com/tscolari/memcached-broker/app"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Catalog app.CfbrokerCatalog `yaml:"catalog"`
}

func Load(filePath string) (Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to open config file: %s\n", err.Error())
	}

	return Parse(data)
}

func Parse(data []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(data, &config)
	return config, err
}
