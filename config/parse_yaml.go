package config

import (
	"launchpad.net/goyaml"
)

type HookDef struct {
	Name          string "name"
	Regexp        string "regexp"
	ShellTemplate string "shell"
	UrlTemplate   string "url"
	FileTemplate  string "file"
}

type Config struct {
	Port     string    "port"
	Host     string    "host"
	User     string    "user"
	HookDefs []HookDef "hooks"
}

func (d *HookDef) GetName() string {
	return d.Name
}

func (d *HookDef) GetRegexp() string {
	return d.Regexp
}

func (d *HookDef) GetShellTemplate() string {
	return d.ShellTemplate
}

func (d *HookDef) GetFileTemplate() string {
	return d.FileTemplate
}

func (d *HookDef) GetUrlTemplate() string {
	return d.UrlTemplate
}

func ParseYaml(yaml []byte) (*Config, error) {
	var config Config
	err := goyaml.Unmarshal(yaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
