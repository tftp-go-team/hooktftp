
package main

import (
	"os"
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


func ParseConfig(path string) (*Config, error) {
	config := &Config{}

	var tmp struct {
		Port string
		Hooks [][]string
	}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
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

