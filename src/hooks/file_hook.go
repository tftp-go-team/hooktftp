package hooks

import (
	"io"
	"os"
	"regexp"

	"github.com/tftp-go-team/libgotftp/src"
)

var pathEscape = regexp.MustCompile("\\.\\.\\/")

var FileHook = HookComponents{
	func(path string, _ tftp.Request) (io.ReadCloser, io.ReadCloser, int, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, nil, -1, err
		}

		// get the file size
		stat, err := file.Stat()
		if err != nil {
			return nil, nil, -1, err
		}
		return file, nil, int(stat.Size()), nil
	},
	func(s string) string {
		return pathEscape.ReplaceAllLiteralString(s, "")
	},
}
