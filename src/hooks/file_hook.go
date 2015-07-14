package hooks

import (
	"io"
	"os"
	"regexp"
)

var pathEscape = regexp.MustCompile("\\.\\.\\/")

var FileHook = HookComponents{
	func(path string) (io.ReadCloser, int, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, -1, err
		}

		// get the file size
		stat, err := file.Stat()
		if err != nil {
			return nil, -1, err
		}
		return file, int(stat.Size()), nil
	},
	func(s string) string {
		return pathEscape.ReplaceAllLiteralString(s, "")
	},
}
