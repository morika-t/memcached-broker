package config

type Instance struct {
	Host           string   `yaml:"host"`
	Port           string   `yaml:"port"`
	OrganizationID string   `yaml:"organization_guid"`
	SpaceID        string   `yaml:"space_guid"`
	ServiceID      string   `yaml:"service_id"`
	PlanID         string   `yaml:"plan_id"`
	Bindings       []string `yaml:"bindings"`
}
