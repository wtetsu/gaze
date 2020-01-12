package config

// Default returns the default configuration
func Default() string {
	return `# Gaze configuration(priority: default < ~/.gaze.yml < ./.gaze.yaml < -f option)
commands:
- ext: .d
  run: dmd -run "{{file}}"
- ext: .js
  run: node "{{file}}"
- ext: .go
  run: go run "{{file}}"
- ext: .groovy
  run: groovy "{{file}}"
- ext: .php
  run: php "{{file}}"
- ext: .pl
  run: perl "{{file}}"
- ext: .py
  run: python "{{file}}"
- ext: .rb
  run: ruby "{{file}}"
- ext: .java
  run: java "{{file}}"
- re: ^Dockerfile$
  run: docker build -f "{{file}}" .
- ext: .c
  run: gcc "{{file}}" && ./a.out
- ext: .cpp
  run: gcc "{{file}}" && ./a.out
- ext: .kts
  run: kotlinc -script "{{file}}"
# - ext: .sh
#   run: sh "{{file}}"
`
}
