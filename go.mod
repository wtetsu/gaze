module github.com/wtetsu/gaze

go 1.13

replace github.com/fsnotify/fsnotify => ./vendor/fsnotify

require (
	github.com/bmatcuk/doublestar v1.2.2
	github.com/cbroglie/mustache v1.0.1
	github.com/fsnotify/fsnotify v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v3 v3.0.0-20191120175047-4206685974f2
)
