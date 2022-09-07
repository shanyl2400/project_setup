package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var config Config

type Config struct {
	Remove []string `yaml:"remove"`
	Move   []string `yaml:"move"`

	Source string `yaml:"src"`
	Backup string `yaml:"backup"`
	Setup  string `yaml:"setup"`

	MoveCpy    string `yaml:"move_cpy"`
	ReleaseCpy string `yaml: "release_cpy"`

	Ignore []string `yaml:"ignore"`
}

func loadConfig(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Read %v failed, err: %v", path, err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Invalid %v format, err: %v", path, err)
	}

	config.Source = fillLastSeparator(config.Source)
	config.Backup = fillLastSeparator(config.Backup)
	config.Setup = fillLastSeparator(config.Setup)
	config.MoveCpy = fillLastSeparator(config.MoveCpy)
	config.ReleaseCpy = fillLastSeparator(config.ReleaseCpy)
}

func fillLastSeparator(path string) string {
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		return path + string(os.PathSeparator)
	}
	return path
}
