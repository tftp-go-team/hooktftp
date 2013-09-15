
package main

import (
	"testing"
	"strings"
)

var jsonBlob = string(`{
   "port": "1234",
   "hooks": [
	   [".*custom.*", "echo jea"]
	]
}`)

func TestParseConfig(t *testing.T) {

	reader := strings.NewReader(jsonBlob)
	config, err := ParseConfig(reader)
	if err != nil {
		t.Fatal("Failed to parse", err)
	}

	if config.Port != "1234" {
		t.Fatal("Bad port", config.Port)
	}

}

