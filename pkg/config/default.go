/**
 * Gaze (https://github.com/wtetsu/gaze/)
 * Copyright 2020-present wtetsu
 * Licensed under MIT
 */

package config

import (
	_ "embed"
)

//go:embed default.yml
var defaultYml string

// Default returns the default configuration
func Default() string {
	return defaultYml
}
