package hooks

import (
	"io"
	"os"
	"regexp"
)

var pathEscape = regexp.MustCompile("\\.\\.\\/")

var FileHook = HookComponents{
	func(path string) (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	},
	func(s string) string {
		return pathEscape.ReplaceAllLiteralString(s, "")
	},
}
