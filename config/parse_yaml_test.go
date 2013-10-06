package config

import (
	"reflect"
	"testing"
)

type yamlExample struct {
	yaml           string
	expectedConfig Config
}

var commandHookTest = yamlExample{
	`
hooks:
  - regexp: ^(.*)$
    command: echo hello
`,
	Config{
		HookDefs: []HookDef{
			{
				Regexp:  "^(.*)$",
				Command: "echo hello",
			},
		},
	},
}

var urlHookTest = yamlExample{
	`
hooks:
  - regexp: ^(.*)$
    url: http://hello
`,
	Config{
		HookDefs: []HookDef{
			{
				Regexp:  "^(.*)$",
				Url: "http://hello",
			},
		},
	},
}

var fileHookTest = yamlExample{
	`
hooks:
  - regexp: ^(.*)$
    file: /var/tftpboot/default
`,
	Config{
		HookDefs: []HookDef{
			{
				Regexp:  "^(.*)$",
				File: "/var/tftpboot/default",
			},
		},
	},
}



var parseTests = []yamlExample{
	{
		"port: 1234",
		Config{Port: "1234"},
	},
	{
		"host: 0.0.0.0",
		Config{Host: "0.0.0.0"},
	},
	{
		"user: hook",
		Config{User: "hook"},
	},
	{
		`
port: 1234
host: 0.0.0.0
user: hook
`,
		Config{
			Port: "1234",
			User: "hook",
			Host: "0.0.0.0",
		},
	},
	commandHookTest,
	urlHookTest,
}

func TestExamples(t *testing.T) {

	for _, example := range parseTests {
		config, err := ParseYaml([]byte(example.yaml))
		if err != nil {
			t.Error("Failed to parse", example.yaml, "because:", err)
			return
		}

		if !reflect.DeepEqual(*config, example.expectedConfig) {
			t.Error(config, "!=", example.expectedConfig)
		}

	}
}

func TestMissingFieldsAreEmpty(t *testing.T) {
	config, err := ParseYaml([]byte("port: 1234"))
	if err != nil {
		t.Error("Failed to parse", err)
		return
	}
	if config.Host != "" {
		t.Error("Host was not empty:", config.Host)
	}
}
