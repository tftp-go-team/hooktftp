package regexptransform

import (
	"errors"
	"testing"

	"github.com/tftp-go-team/hooktftp/internal/config"
)

func dummyEscape(s string, extraArgs config.HookExtraArgs) (string, error) {
	return s, nil
}

var INVALID_ESCAPE = errors.New("escape failed")

func erroneousEscape(s string, extraArgs config.HookExtraArgs) (string, error) {
	return "", INVALID_ESCAPE
}

type TransformTest struct {
	regexp   string
	template string
	escape   Escape
	input    string
	output   string
	err      error
}

var TransformTests = []TransformTest{
	{
		".*",
		"$0",
		dummyEscape,
		"/var/tftpboot/hello",
		"/var/tftpboot/hello",
		nil,
	},
	{
		".*",
		"$0",
		erroneousEscape,
		"/var/tftpboot/hello",
		"",
		INVALID_ESCAPE,
	},
	{
		"/var/tftpboot/(.*)$",
		"http://localhost/get/$1",
		dummyEscape,
		"/var/tftpboot/hello",
		"http://localhost/get/hello",
		nil,
	},
	{
		"/var/tftpboot/(.*)$",
		"http://localhost/full$0",
		dummyEscape,
		"/var/tftpboot/hello",
		"http://localhost/full/var/tftpboot/hello",
		nil,
	},
	{
		"/var/tftpboot/(.*)$",
		"http://localhost/get$0",
		dummyEscape,
		"nomatch",
		"",
		NO_MATCH,
	},
	{
		"/var/tftpboot/(.*)$",
		"http://localhost/get/$1",
		func(s string, extraArgs config.HookExtraArgs) (string, error) {
			return "'" + s + "'", nil
		},
		"/var/tftpboot/hello",
		"http://localhost/get/'hello'",
		nil,
	},
	{
		"/var/tftpboot/(.+)/(.+).txt$",
		"http://localhost/get/$1/$2",
		dummyEscape,
		"/var/tftpboot/foo/bar.txt",
		"http://localhost/get/foo/bar",
		nil,
	},
	{
		"/var/tftpboot/(.+)/(.+).txt$",
		"http://localhost/get/$1",
		dummyEscape,
		"/var/tftpboot/foo/bar.txt",
		"http://localhost/get/foo",
		nil,
	},
	{
		"/var/tftpboot/(.+)$",
		"http://localhost/get/$1/$2",
		dummyEscape,
		"/var/tftpboot/foo/bar.txt",
		"",
		BAD_GROUPS,
	},
}

func TestTranforms(t *testing.T) {
	for _, tt := range TransformTests {

		transform, err := NewRegexpTransform(
			tt.regexp,
			tt.template,
			tt.escape,
			nil,
		)
		if err != nil {
			t.Error("failed to compile transform:", err, tt)
			return
		}

		res, err := transform(tt.input)
		if err != tt.err {
			t.Errorf("Unexpected transform error '%v' expected '%v'\nfor '%v'", err, tt.err, tt)
			return
		}
		if res != tt.output {
			t.Errorf("Unexpected output '%v', expected '%v'\nfor '%v'", res, tt.output, tt)
		}

	}
}
