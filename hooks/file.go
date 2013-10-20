
package hooks

import (
	"os"
	"io"
)

var FileHook = HookComponents{
	func (path string) (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		return file, nil
	},
	func (s string) string{
		return s
	},
}
