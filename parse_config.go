
package main

import (
	"os"
	"io"
	"encoding/json"
	"regexp"
)

type Hook struct {
	Regexp *regexp.Regexp
	Command string
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

