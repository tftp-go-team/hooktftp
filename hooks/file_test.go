package hooks

import (
	"github.com/epeli/dyntftp/config"
	"testing"
)

type hookTestCase struct {
	regexp   string
	template string
	input    string
}

var hookTestCases = []hookTestCase{
	{
		".*",
		"file_fixture.txt",
		"anything",
	},
	{
		"^extension:(.*)$",
		"file_fixture.$1",
		"extension:txt",
	},
}

func TestFile(t *testing.T) {
	for _, testCase := range hookTestCases {
		hookDef := &config.HookDef{
			Regexp:       testCase.regexp,
			FileTemplate: testCase.template,
		}

		hook, err := CompileHook(hookDef)
		if err != nil {
			t.Error("Failed to compile", hookDef, err)
			return
		}

		file, err := hook(testCase.input)
		if err != nil {
			t.Error("Failed to execute hook", hookDef, err)
			return
		}

		data := make([]byte, 5)
		_, err = file.Read(data)
		if err != nil {
			t.Error("Failed to read file", file)
			return
		}

		if string(data) != "hello" {
			t.Error("Expected to find hello from file but got:", string(data))
		}
	}

}
