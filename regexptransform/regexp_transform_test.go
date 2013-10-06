package regexptransform

import (
	"testing"
)

func dummyEscape(s string) string {
	return s
}

type TransformTest struct {
	regexp string
	template string
	escape Escape
	input string
	output string
	err error
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
		func(s string) string {
			return "'" + s + "'"
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


