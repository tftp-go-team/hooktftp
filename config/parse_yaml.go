package config

import (
	"gopkg.in/yaml.v1"
)

type HookDef struct {
	Description string "description"
	Type        string "type"
	Regexp      string "regexp"
	Template    string "template"
}

type Config struct {
	Port     string    "port"
	Host     string    "host"
	User     string    "user"
	HookDefs []HookDef "hooks"
}

func (d *HookDef) GetType() string {
	return d.Type
}

func (d *HookDef) GetTemplate() string {
	return d.Template
}

func (d *HookDef) GetDescription() string {
	return d.Description
}

func (d *HookDef) GetRegexp() string {
	return d.Regexp
}

func ParseYaml(yamlv1 []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(yamlv1, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
