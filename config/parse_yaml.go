package config

import (
	"launchpad.net/goyaml"
)

type HookDef struct {
	Regexp          string "regexp"
	CommandTemplate string "command"
	UrlTemplate     string "url"
	FileTemplate    string "file"
}

type Config struct {
	Port     string    "port"
	Host     string    "host"
	User     string    "user"
	HookDefs []HookDef "hooks"
}

func ParseYaml(yaml []byte) (*Config, error) {
	var config Config
	err := goyaml.Unmarshal(yaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
