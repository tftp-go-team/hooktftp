package config

import (
	"reflect"
	"testing"
)

type yamlExample struct {
	yaml           string
	expectedConfig Config
}

var shellHookTest = yamlExample{
	`
hooks:
  - type: shell
    regexp: ^(.*)$
    template: echo hello
`,
	Config{
		HookDefs: []HookDef{
			{
				Type:     "shell",
				Regexp:   "^(.*)$",
				Template: "echo hello",
			},
		},
	},
}

var shellHookWithWhitelistTest = yamlExample{
	`
hooks:
  - type: shell
    regexp: ^(.*)$
    template: echo hello
    whitelist:
      - allowed
      - allowed_too
`,
	Config{
		HookDefs: []HookDef{
			{
				Type:      "shell",
				Regexp:    "^(.*)$",
				Template:  "echo hello",
				Whitelist: []string{"allowed", "allowed_too"},
			},
		},
	},
}

var httpHookTest = yamlExample{
	`
hooks:
  - type: http
    regexp: ^(.*)$
    template: http://hello
`,
	Config{
		HookDefs: []HookDef{
			{
				Type:     "http",
				Regexp:   "^(.*)$",
				Template: "http://hello",
			},
		},
	},
}

var fileHookTest = yamlExample{
	`
hooks:
  - type: file
    regexp: ^(.*)$
    template: /var/tftpboot/default
`,
	Config{
		HookDefs: []HookDef{
			{
				Regexp:   "^(.*)$",
				Template: "/var/tftpboot/default",
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
	shellHookTest,
	shellHookWithWhitelistTest,
	httpHookTest,
}

func TestExamples(t *testing.T) {

	for _, example := range parseTests {
		config, err := ParseYaml([]byte(example.yaml))
		if err != nil {
			t.Error("Failed to parse", example.yaml, "because:", err)
			return
		}

		if !reflect.DeepEqual(*config, example.expectedConfig) {
			t.Errorf("Expected '%v' to match '%v'", config, example.expectedConfig)
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
