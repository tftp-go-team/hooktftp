package hooks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"config"
	"io/ioutil"
)

type hookTestCase struct {
	hookDef        *config.HookDef
	input          string
	expected       string
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
			w.WriteHeader(200)
			fmt.Fprintln(w, "RES:web.txt")
			return
		}

	}))
	defer ts.Close()

	var hookTestCases = []hookTestCase{
		{
			&config.HookDef{
				Type:     "file",
				Regexp:   ".*",
				Template: "file_fixture.txt",
			},
			"anything",
			"filecontents",
			noError,
		},
		{
			&config.HookDef{
				Type:     "file",
				Regexp:   ".*",
				Template: "$0",
			},
			"../file_fixture.txt",
			"filecontents",
			noError,
		},
		{
			&config.HookDef{
				Type:     "file",
				Regexp:   "^extension:(.*)$",
				Template: "file_fixture.$1",
			},
			"extension:txt",
			"filecontents",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   ".*",
				Template: "echo shellhello",
			},
			"anything",
			"shellhello",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   "foo:(.*)",
				Template: "echo $1",
			},
			"foo:haha",
			"haha",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   "foo:(.*)",
				Template: "echo $1",
			},
			"foo:$(hostname)",
			"$(hostname)",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   "foo:(.*)",
				Template: "echo $1",
			},
			"foo:`hostname`",
			"`hostname`",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   "foo:(.*)",
				Template: "echo $1",
			},
			"foo:$HOME",
			"$HOME",
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   "(.*)",
				Template: "echo $1; exit 1",
			},
			"foo",
			"foo",
			noError,
		},
		{
			&config.HookDef{
				Type:     "http",
				Regexp:   "url\\/(.+)$",
				Template: ts.URL + "/test/$1",
			},
			"url/web.txt",
			"RES:web.txt",
			noError,
		},
		{
			&config.HookDef{
				Type:     "http",
				Regexp:   "url\\/(.+)$",
				Template: ts.URL + "/bad",
			},
			"url/bad.txt",
			"bad response",
			func(err error) error {
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

		file, _, err := hook(testCase.input)
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

		data, err := ioutil.ReadAll(file)
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
