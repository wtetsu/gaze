package config

import (
	"errors"
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

// LoadConfig loads a configuration file.
// Priority: default < ~/.gaze.yml < ./.gaze.yaml < -f option)
func LoadConfig() (*Config, error) {
	configPath, err := searchConfigPath()
	logger.Debug("configPath: " + configPath)

	var bytes []byte
	if err == nil {
		bytes, err = ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
	} else {
		bytes = []byte(Default())
	}

	entries, err := parseConfig(bytes)
	if err != nil {
		return nil, err
	}

	var config Config
	config.Commands = *entries

	return prepare(&config), nil
}

func prepare(configs *Config) *Config {
	for i := 0; i < len(configs.Commands); i++ {
		reStr := configs.Commands[i].Re
		re, err := regexp.Compile(reStr)
		if err == nil {
			configs.Commands[i].re = re
		}
	}
	return configs
}

func searchConfigPath() (string, error) {
	const CONFIG = ".gaze.yml"

	filepath.FromSlash("path string")
	filepath.ToSlash("path string")

	path1 := "./" + CONFIG
	if fs.IsFile(path1) {
		return path1, nil
	}

	home := homeDirPath()
	if home != "" {
		path2 := path.Join(home, CONFIG)
		if fs.IsFile(path2) {
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
	if c.Ext == "" && c.re == nil {
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

	return false
}
