package config

type State struct {
	Capacity  int        `yaml:"capacity"`
	Instances []Instance `yaml:"instances"`
}

type Instance struct {
	ID               string `yaml:"id"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	OrganizationGUID string `yaml:"organization_guid"`
	SpaceGUID        string `yaml:"space_guid"`
}
