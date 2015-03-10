package hooks

import (
	"fmt"
	"io"
	"net/http"
)

var HTTPHook = HookComponents{
	func(url string) (io.ReadCloser, error) {

		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return nil, fmt.Errorf("Bad response '%v' from %v", res.Status, url)
		}

		return res.Body, nil

	},
	func(s string) string {
		return s
	},
}
