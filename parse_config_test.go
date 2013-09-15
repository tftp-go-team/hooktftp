package main

import (
	"strings"
	"testing"
)

var jsonBlob = string(`{
   "port": "1234",
   "hooks": [
   [".*custom.*", "sh:echo jea"]
	]
}`)

func createTestConfig() (*Config, error) {
	reader := strings.NewReader(jsonBlob)
	return ParseConfig(reader)
}

func TestParseConfigPort(t *testing.T) {

	config, err := createTestConfig()
	if err != nil {
		t.Fatal("Failed to parse config", err)
		return
	}

	if config.Port != "1234" {
		t.Fatal("Bad port", config.Port)
		return
	}
}
func TestParseConfigShHook(t *testing.T) {

	config, err := createTestConfig()
	if err != nil {
		t.Fatal("Failed to parse config", err)
		return
	}

	if match := config.Hooks[0].Regexp.Find([]byte("foocustombar")); match == nil {
		t.Fatal("regexp matcher failed")
		return
	}

	out, err := config.Hooks[0].Execute()
	if err != nil {
		t.Fatal("hook execute err", err.Error())
		return
	}

	b := make([]byte, 3)
	if _, err = out.Read(b); err != nil {
		t.Fatal("Read fail", err)
		return
	}

	if string(b) != "jea" {
		t.Fatal("bad output from execute", string(b))
	}


}
