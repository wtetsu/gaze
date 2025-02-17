/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	"errors"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"github.com/cbroglie/mustache"
	"github.com/wtetsu/gaze/pkg/gutil"
	"github.com/wtetsu/gaze/pkg/logger"
	"gopkg.in/yaml.v3"
)

// For deserialize
type rawConfig struct {
	Commands []rawCommand
	Log      *rawLog
}

// For deserialize
type rawCommand struct {
	Ext string
	Cmd string
	Re  string
}

// For deserialize
type rawLog struct {
	Start string
	End   string
}

// Config represents Gaze configuration
type Config struct {
	Commands []Command
	Log      *Log
}

// Command represents Gaze configuration
type Command struct {
	Ext string
	Cmd string
	re  *regexp.Regexp
}

type Log struct {
	start *mustache.Template
	end   *mustache.Template
}

// New returns a new Config.
func NewWithFixedCommand(command string) (*Config, error) {
	if command == "" {
		return nil, errors.New("empty command")
	}
	fixedCommand := rawCommand{Cmd: command, Re: "."}
	loadedRawConfig, err := loadPreferredRawConfig(homeDirPath())
	if err != nil {
		return nil, err
	}

	config := rawConfig{Commands: []rawCommand{fixedCommand}, Log: loadedRawConfig.Log}
	return toConfig(&config), nil
}

// LoadPreferredConfig loads a configuration file.
// Priority: default < ~/.gaze.yml < ~/.config/gaze/gaze.yml < -f option)
func LoadPreferredConfig() (*Config, error) {
	rawConfig, err := loadPreferredRawConfig(homeDirPath())
	if err != nil {
		return nil, err
	}
	return toConfig(rawConfig), nil
}

func loadPreferredRawConfig(home string) (*rawConfig, error) {
	configPath := searchConfigPath(home)

	if configPath != "" {
		logger.Info("config: " + configPath)
		parsedRawConfig, err := parseRawConfigFromFile(configPath)
		if err != nil {
			return nil, err
		}
		return parsedRawConfig, nil
	}

	logger.Info("config: (default)")
	return defaultRawConfig(), nil
}

func defaultRawConfig() *rawConfig {
	bytes := []byte(Default())
	rawConfig, _ := parseRawConfigFromBytes(bytes) // Parse error never occurs
	return rawConfig
}

// LoadConfigFromFile loads a configuration file.
func LoadConfigFromFile(configPath string) (*Config, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	rawConfig, err := parseRawConfigFromBytes(bytes)
	if err != nil {
		return nil, err
	}

	return toConfig(rawConfig), nil
}

func toConfig(rawConfig *rawConfig) *Config {
	resultConfig := &Config{}
	if len(rawConfig.Commands) == 0 {
		logger.Notice("No commands defined in the configuration file. Gaze will not function properly.")
	}

	for i := 0; i < len(rawConfig.Commands); i++ {
		rawCmd := &rawConfig.Commands[i]

		if rawCmd.Cmd == "" {
			logger.Error("Empty cmd (%d)", i)
			continue
		}
		if rawCmd.Ext == "" && rawCmd.Re == "" {
			logger.Debug("Both ext and re are empty (%d)", i)
			continue
		}

		if rawCmd.Re != "" {
			re, err := regexp.Compile(rawCmd.Re)
			if err == nil {
				resultConfig.Commands = append(resultConfig.Commands, Command{Cmd: rawCmd.Cmd, Ext: rawCmd.Ext, re: re})
			} else {
				logger.Error("Failed to compile regexp: " + err.Error())
			}
			continue
		}

		if rawCmd.Ext != "" {
			resultConfig.Commands = append(resultConfig.Commands, Command{Cmd: rawCmd.Cmd, Ext: rawCmd.Ext})
			continue
		}
	}

	sourceLog := rawConfig.Log
	if sourceLog == nil {
		defaultRawConfig, _ := parseRawConfigFromBytes([]byte(Default())) // Parse error never occurs
		sourceLog = defaultRawConfig.Log
	}

	start := parseMustacheTemplate(sourceLog.Start)
	end := parseMustacheTemplate(sourceLog.End)

	resultConfig.Log = &Log{start: start, end: end}

	return resultConfig
}

// parseMustacheTemplate parses a mustache template that tolerates errors
func parseMustacheTemplate(source string) *mustache.Template {
	template, err := mustache.ParseStringRaw(source, true)
	if err != nil {
		logger.Error("Failed to parse template: %s: %s", err.Error(), source)
		return nil
	}
	return template
}

func searchConfigPath(home string) string {
	if !gutil.IsDir(home) {
		return ""
	}
	configDir := path.Join(home, ".config", "gaze")
	for _, n := range []string{"gaze.yml", "gaze.yaml"} {
		candidate := path.Join(configDir, n)
		if gutil.IsFile(candidate) {
			return candidate
		}
	}
	for _, n := range []string{".gaze.yml", ".gaze.yaml"} {
		candidate := path.Join(home, n)
		if gutil.IsFile(candidate) {
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

func parseRawConfigFromFile(path string) (*rawConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseRawConfigFromBytes(bytes)
}

func parseRawConfigFromBytes(fileBuffer []byte) (*rawConfig, error) {
	rawConfig := rawConfig{}
	err := yaml.Unmarshal(fileBuffer, &rawConfig)
	if err != nil {
		return nil, err
	}

	return &rawConfig, nil
}

// Match return true is filePath meets the condition
func (c *Command) Match(filePath string) bool {
	if filePath == "" {
		return false
	}

	if c.Ext != "" && c.re == nil {
		return c.Ext == filepath.Ext(filePath)
	}
	if c.Ext == "" && c.re != nil {
		return c.re.MatchString(filePath)
	}

	// Both are set
	return c.Ext == filepath.Ext(filePath) && c.re.MatchString(filePath)
}

func (l *Log) RenderStart(params map[string]string) string {
	return renderLog(l.start, params)
}

func (l *Log) RenderEnd(params map[string]string) string {
	return renderLog(l.end, params)
}

func renderLog(tmpl *mustache.Template, params map[string]string) string {
	if tmpl == nil {
		return ""
	}
	log, err := tmpl.Render(params)
	if err != nil {
		logger.Error("Failed to render log: %s", err)
		return ""
	}
	return log
}
