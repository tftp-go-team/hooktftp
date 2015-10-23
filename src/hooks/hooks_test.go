package hooks

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tftp-go-team/hooktftp/src/config"
	"github.com/tftp-go-team/libgotftp/src"
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

	clientAddr := net.Addr(&net.TCPAddr{
		IP:   net.ParseIP("198.51.100.13"),
		Port: 63233,
	})

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
				Type:     "shell",
				Regexp:   ".*",
				Template: "sh -c 'echo $CLIENT_ADDR'",
			},
			"anything",
			clientAddr.String(),
			noError,
		},
		{
			&config.HookDef{
				Type:     "shell",
				Regexp:   ".*",
				Template: "",
			},
			"anything",
			"anything",
			func(err error) error {
				if err == nil {
					return fmt.Errorf("Bad url response test failed: Expected to have an error")
				}
				return nil
			},
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

		fakeRequest := tftp.Request{Addr: &clientAddr}

		file, _, _, err := hook(testCase.input, fakeRequest)
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
