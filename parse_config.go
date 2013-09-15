
package main

import (
	"fmt"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Hook struct {
	Regexp *regexp.Regexp
	Command string
}

func (h *Hook) Execute() (io.Reader, error) {
	if strings.HasPrefix(h.Command, "sh:") {
		return h.ShExecute()
	}
	// TODO: implement http fetch

	return nil, fmt.Errorf("Unknown command type %v", h.Command)
}

func (h *Hook) ShExecute() (io.Reader, error) {
	cmd := exec.Command("sh", "-c", h.Command[3:])
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	return reader, err
}


type Config struct {
	Port string
	Hooks []*Hook
}

func ParseConfigFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseConfig(file)
}


func ParseConfig(reader io.Reader) (*Config, error) {
	config := &Config{}

	var tmp struct {
		Port string
		Hooks [][]string
	}

	jsonDecoder := json.NewDecoder(reader)
	if err := jsonDecoder.Decode(&tmp); err != nil {
		return config, err
	}

	config.Port = tmp.Port
	for _, rawhook := range tmp.Hooks {
		r, err := regexp.Compile(rawhook[0])
		if err != nil {
			return config, err
		}

		config.Hooks = append(config.Hooks, &Hook{
			r,
			rawhook[1],
		})
	}

	return config, nil

}

