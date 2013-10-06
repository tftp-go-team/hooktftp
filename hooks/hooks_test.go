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
	{
		&config.HookDef{
			Regexp:       ".*",
			ShellTemplate: "echo shellhello",
		},
		"anything",
		"shellhello",
	},
	{
		&config.HookDef{
			Regexp:       "foo:(.*)",
			ShellTemplate: "echo $1",
		},
		"foo:haha",
		"haha",
	},
	{
		&config.HookDef{
			Regexp:       "foo:(.*)",
			ShellTemplate: "echo $1",
		},
		"foo:$(hostname)",
		"$(hostname)",
	},
	{
		&config.HookDef{
			Regexp:       "foo:(.*)",
			ShellTemplate: "echo $1",
		},
		"foo:`hostname`",
		"`hostname`",
	},
	{
		&config.HookDef{
			Regexp:       "foo:(.*)",
			ShellTemplate: "echo $1",
		},
		"foo:$HOME",
		"$HOME",
	},
	{
		&config.HookDef{
			Regexp:       "(.*)",
			ShellTemplate: "echo $1; exit 1",
		},
		"foo",
		"foo",
	},
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

		data := make([]byte, 20)
		_, err = file.Read(data)
		if err != nil {
			t.Error("Failed to read file", testCase.hookDef,  file)
			return
		}

		res := string(data[:len(testCase.expected)])

		if res != testCase.expected {
			t.Errorf("Expected to find '%v' from file but got '%v'", testCase.expected,  res)
		}
	}

}
