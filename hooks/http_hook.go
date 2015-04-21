package hooks

import (
	"fmt"
	"io"
	"net/http"
)

var HTTPHook = HookComponents{
	func(url string) (io.ReadCloser, error) {
		client := &http.Client {
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) == 0 {
					return nil
				}

				for key, val := range via[0].Header {
					if key != "Referer" {
						req.Header[key] = val
					}
				}

				return nil
			},
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "hooktftp v0.9.1")

		res, err := client.Do(req)
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
