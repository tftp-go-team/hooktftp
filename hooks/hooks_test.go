package hooks

import (
	"fmt"
	"github.com/epeli/hooktftp/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

type hookTestCase struct {
	hookDef  *config.HookDef
	input    string
	expected string
	errorValidator func(error) error
}

func noError(err error) error {
	return err
}

func TestHooks(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/bad" {
			w.WriteHeader(500)
			fmt.Fprintln(w, "bad response")
		}

		if r.URL.String() == "/test/web.txt" {
			fmt.Fprintln(w, "RES:web.txt")
			return
		}

	}))
	defer ts.Close()

	var hookTestCases = []hookTestCase{
		{
			&config.HookDef{
				Regexp:       ".*",
				FileTemplate: "file_fixture.txt",
			},
			"anything",
			"filecontents",
			noError,
		},
		{
			&config.HookDef{
				Regexp:       "^extension:(.*)$",
				FileTemplate: "file_fixture.$1",
			},
			"extension:txt",
			"filecontents",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        ".*",
				ShellTemplate: "echo shellhello",
			},
			"anything",
			"shellhello",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        "foo:(.*)",
				ShellTemplate: "echo $1",
			},
			"foo:haha",
			"haha",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        "foo:(.*)",
				ShellTemplate: "echo $1",
			},
			"foo:$(hostname)",
			"$(hostname)",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        "foo:(.*)",
				ShellTemplate: "echo $1",
			},
			"foo:`hostname`",
			"`hostname`",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        "foo:(.*)",
				ShellTemplate: "echo $1",
			},
			"foo:$HOME",
			"$HOME",
			noError,
		},
		{
			&config.HookDef{
				Regexp:        "(.*)",
				ShellTemplate: "echo $1; exit 1",
			},
			"foo",
			"foo",
			noError,
		},
		{
			&config.HookDef{
				Regexp:       "url\\/(.+)$",
				UrlTemplate: ts.URL + "/test/$1",
			},
			"url/web.txt",
			"RES:web.txt",
			noError,
		},
		{
			&config.HookDef{
				Regexp:       "url\\/(.+)$",
				UrlTemplate: ts.URL + "/bad",
			},
			"url/bad.txt",
			"bad response",
			func (err error) error {
				if err == nil {
					return fmt.Errorf("Bad url response test failed: Expected to have an error")
				}
				return nil
			},
		},
	}

	for _, testCase := range hookTestCases {
		hook, err := CompileHook(testCase.hookDef)
		if err != nil {
			t.Error("Failed to compile", testCase.hookDef, err)
			return
		}

		file, err := hook(testCase.input)
		if err == NO_MATCH {
			t.Error(testCase.hookDef.Regexp, "does not match with", testCase.input)
		}

		if err = testCase.errorValidator(err); err != nil {
			t.Error("Failed to execute hook:", testCase.hookDef, err)
			return
		}

		if file == nil {
			return
		}

		data := make([]byte, 20)
		_, err = file.Read(data)
		if err != nil {
			t.Error("Failed to read file", testCase.hookDef, file)
			return
		}

		res := string(data[:len(testCase.expected)])

		if res != testCase.expected {
			t.Errorf("Expected to find '%v' from file but got '%v'", testCase.expected, res)
		}
	}
}
