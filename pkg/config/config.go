package config

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"

	"github.com/wtetsu/gaze/pkg/file"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads a configuration file.
// Priority: ./.gaze.yaml > ~/.gaze.yml > default
func LoadConfig() ([]Config, error) {
	configPath, err := searchConfigPath()

	var bytes []byte
	if err == nil {
		bytes, err = ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	} else {
		bytes = []byte(Default())
	}

	config, err := parseConfig(bytes)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func searchConfigPath() (string, error) {
	const CONFIG = ".gaze.yml"

	filepath.FromSlash("path string")
	filepath.ToSlash("path string")

	path1 := "./" + CONFIG
	if file.Exist(path1) {
		return path1, nil
	}

	home := homeDirPath()
	if home != "" {
		path2 := path.Join(home, CONFIG)
		if file.Exist(path2) {
			return path2, nil
		}
	}

	return "", errors.New("Config file not found")
}

func homeDirPath() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.ToSlash(currentUser.HomeDir)
}

func parseConfig(fileBuffer []byte) ([]Config, error) {
	data := make([]Config, 20)
	err := yaml.Unmarshal(fileBuffer, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Config represents Gaze configuration
type Config struct {
	Ext   string
	Run   string
	Match string
}

// Default returns the default configuration
func Default() string {
	return `# Gaze configuration(priority: ./.gaze.yaml > ~/.gaze.yml > default)
- ext: .d
  run: dmd -run {{file}}
- ext: .js
  run: node {{file}}
- ext: .go
  run: go run {{file}}
- ext: .groovy
  run: groovy {{file}}
- ext: .php
  run: php {{file}}
- ext: .pl
  run: perl {{file}}
- ext: .py
  run: python {{file}}
- ext: .rb
  run: ruby {{file}}
- ext: .java
  run: java {{file}}
- match: ^Dockerfile$
  run: docker build -f {{file}} .
- ext: .c
  run: gcc {{file}} && ./a.out
- ext: .cpp
  run: gcc {{file}} && ./a.out
- ext: .kts
  run: kotlinc -script {{file}}	
# - ext: .sh
#   run: sh {{file}}
`
}
