/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	"io/ioutil"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"github.com/wtetsu/gaze/pkg/fs"
	"github.com/wtetsu/gaze/pkg/logger"
	"gopkg.in/yaml.v3"
)

// Config represents Gaze configuration
type Config struct {
	Commands []Command
}

// Command represents Gaze configuration
type Command struct {
	Ext string
	Run string
	Re  string
	re  *regexp.Regexp
}

// New returns
func New(command string) *Config {
	fixedCommand := Command{Re: ".", Run: command}
	config := Config{Commands: []Command{fixedCommand}}
	return prepare(&config)
}

// InitConfig loads a configuration file.
// Priority: default < ~/.gaze.yml < ./.gaze.yaml < -f option)
func InitConfig(fileNameList []string) (*Config, error) {
	configPath := searchConfigPath(fileNameList)
	return makeConfig(configPath)
}

func makeConfig(configPath string) (*Config, error) {
	if configPath != "" {
		logger.Info("config: " + configPath)
		return LoadConfig(configPath)
	}

	logger.Info("config: (default)")
	return defaultConfig()
}

func defaultConfig() (*Config, error) {
	bytes := []byte(Default())
	entries, err := parseConfig(bytes)
	if err != nil {
		return nil, err
	}

	config := Config{Commands: *entries}
	return prepare(&config), nil
}

// LoadConfig loads a configuration file.
func LoadConfig(configPath string) (*Config, error) {
	logger.Info("config: " + configPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return makeConfigFromBytes(bytes)
}

func makeConfigFromBytes(bytes []byte) (*Config, error) {
	commands, err := parseConfig(bytes)
	if err != nil {
		return nil, err
	}

	config := Config{Commands: *commands}
	return prepare(&config), nil
}

func prepare(configs *Config) *Config {
	for i := 0; i < len(configs.Commands); i++ {
		reStr := configs.Commands[i].Re
		if reStr == "" {
			continue
		}
		re, err := regexp.Compile(reStr)
		if err == nil {
			configs.Commands[i].re = re
		}
	}
	return configs
}

func searchConfigPath(fileNameList []string) string {
	candidates := []string{}
	for _, n := range fileNameList {
		candidates = append(candidates, "./"+n)
	}

	home := homeDirPath()
	if home != "" {
		for _, n := range fileNameList {
			candidates = append(candidates, path.Join(home, n))
		}
	}
	for _, c := range candidates {
		if fs.IsFile(c) {
			return c
		}
	}
	return ""
}

func homeDirPath() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}
	return filepath.ToSlash(currentUser.HomeDir)
}

func parseConfig(fileBuffer []byte) (*[]Command, error) {
	config := Config{}
	err := yaml.Unmarshal(fileBuffer, &config)
	if err != nil {
		return nil, err
	}
	return &config.Commands, nil
}

// Match return true is filePath meets the condition
func (c *Command) Match(filePath string) bool {
	if filePath == "" {
		return false
	}

	if c.Ext != "" && c.re != nil {
		return c.Ext == filepath.Ext(filePath) && c.re.MatchString(filePath)
	}
	if c.Ext != "" && c.re == nil {
		return c.Ext == filepath.Ext(filePath)
	}
	if c.Ext == "" && c.re != nil {
		return c.re.MatchString(filePath)
	}

	// c.Ext == "" && c.re == nil
	return false
}
