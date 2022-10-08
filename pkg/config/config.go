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
	Cmd string
	Re  string
	re  *regexp.Regexp
}

// New returns a new Config.
func New(command string) *Config {
	fixedCommand := Command{Re: ".", Cmd: command}
	config := Config{Commands: []Command{fixedCommand}}
	return prepare(&config)
}

// InitConfig loads a configuration file.
// Priority: default < ~/.gaze.yml < ~/.config/gaze.yaml < -f option)
func InitConfig() (*Config, error) {
	home := homeDirPath()
	configPath := searchConfigPath(home)
	return makeConfigFromFile(configPath)
}

func makeConfigFromFile(configPath string) (*Config, error) {
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

func searchConfigPath(home string) string {
	if !fs.IsDir(home) {
		return ""
	}
	configDir := path.Join(home, ".config", "gaze")
	for _, n := range []string{"gaze.yml", "gaze.yaml"} {
		candidate := path.Join(configDir, n)
		if fs.IsFile(candidate) {
			return candidate
		}
	}
	for _, n := range []string{".gaze.yml", ".gaze.yaml"} {
		candidate := path.Join(home, n)
		if fs.IsFile(candidate) {
			return candidate
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
