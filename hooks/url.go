package hooks

import (
	"io"
	"net/http"
)

var UrlHook = HookComponents{
	func(url string) (io.ReadCloser, error) {

		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		return res.Body, nil

	},
	func(s string) string {
		return s
	},
}
