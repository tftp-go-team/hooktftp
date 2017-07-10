package config

import (
	"errors"
	"launchpad.net/goyaml"
)

type HookDef struct {
	Description string   "description"
	Type        string   "type"
	Regexp      string   "regexp"
	Template    string   "template"
	Whitelist   []string "whitelist"
	UrlDecode   bool     "urldecode"
}

type Config struct {
	Port     string    "port"
	Host     string    "host"
	User     string    "user"
	HookDefs []HookDef "hooks"
}

type HookExtraArgs map[string]interface{}

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

func (d *HookDef) GetWhitelist() []string {
	return d.Whitelist
}

func (d *HookDef) GetExtraArgs() HookExtraArgs {
	ret := make(HookExtraArgs)
	ret["urldecode"] = d.UrlDecode
	return ret
}

func ParseYaml(yaml []byte) (*Config, error) {
	var config Config
	err := goyaml.Unmarshal(yaml, &config)
	if err != nil {
		return nil, err
	}

	for _, hookdef := range config.HookDefs {
		if hookdef.UrlDecode && hookdef.Type != "http" {
			return nil, errors.New("urldecode option is only valid for the http hook")
		}
	}

	return &config, nil
}
