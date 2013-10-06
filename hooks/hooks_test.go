package hooks

import (
	"github.com/epeli/dyntftp/config"
	"testing"
)

type hookTestCase struct {
	hookDef  *config.HookDef
	input    string
	expected string
}

var hookTestCases = []hookTestCase{
	{
		&config.HookDef{
			Regexp:       ".*",
			FileTemplate: "file_fixture.txt",
		},
		"anything",
		"filecontents",
	},
	{
		&config.HookDef{
			Regexp:       "^extension:(.*)$",
			FileTemplate: "file_fixture.$1",
		},
		"extension:txt",
		"filecontents",
	},
	// {
	// 	&config.HookDef{
	// 		Regexp:       ".*",
	// 		ShellTemplate: "echo shellhello",
	// 	},
	// 	"anything",
	// 	"shellhello",
	// },
}

func TestFile(t *testing.T) {
	for _, testCase := range hookTestCases {
		hook, err := CompileHook(testCase.hookDef)
		if err != nil {
			t.Error("Failed to compile", testCase.hookDef, err)
			return
		}

		file, err := hook(testCase.input)
		if err != nil {
			t.Error("Failed to execute hook", testCase.hookDef, err)
			return
		}

		data := make([]byte, 12)
		_, err = file.Read(data)
		if err != nil {
			t.Error("Failed to read file", file)
			return
		}

		res := string(data)

		if res != testCase.expected {
			t.Error("Expected to find", testCase.expected, "from file but got:", res)
		}
	}

}
