**Gaze is gazing at you**

![GAZE](https://user-images.githubusercontent.com/515948/71816598-828a9700-30c6-11ea-92c8-ca0154e98794.png)

[![Build Status](https://travis-ci.com/wtetsu/gaze.svg?branch=master)](https://travis-ci.com/wtetsu/gaze) [![Go Report Card](https://goreportcard.com/badge/github.com/wtetsu/gaze)](https://goreportcard.com/report/github.com/wtetsu/gaze) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/ec1ab9cfb5b04feba674c1c1440ffb99)](https://www.codacy.com/manual/wtetsu/gaze?utm_source=github.com&utm_medium=referral&utm_content=wtetsu/gaze&utm_campaign=Badge_Grade) [![codecov](https://codecov.io/gh/wtetsu/gaze/branch/master/graph/badge.svg)](https://codecov.io/gh/wtetsu/gaze)

# Gaze

## What is Gaze?

Gaze runs a command, **right after** you save a file.

It greatly helps you focus on writing code!
![gaze02](https://user-images.githubusercontent.com/515948/73607575-1fbfe900-45fb-11ea-813e-6be6bf9ece6d.gif)

## Features:

- Easy to use, out-of-the-box
- React super quickly to your file modifications
- Language-agnostic
- Flexible configuration
- Useful advanced options
  - timeout(useful if you sometimes write infinite loops)
  - restart(useful for server applications)
- Multiplatform (macOS, Windows, Linux)

## Use cases:

🚀Gaze runs a script, **right after** you save it(e.g. Python, Ruby),

You can also use Gaze for these purposes:

- 🚀Gaze runs tests, **right after** you save a Ruby script
- 🚀Gaze runs linter, **right after** you save a JavaScript file
- 🚀Gaze runs "docker build .", **right after** you save Dockerfile
- And so forth...

---

Software development sometimes forces us to execute the same command again and again, by hands!

Let's say, you started writing a really really really simple Python script. You created a.py, wrote 5 lines of code and run "python a.py".
Since the result was not perfect, you edited a.py again, and run "python a.py" again.

Again and again...

Before you realized, you've saved the same files and executes the same command thousands of times!

Gaze runs a command Instead of you, **right after** you edit files.

# Installation

## Brew (Only for OSX)

```
brew tap wtetsu/gaze
brew install gaze
```

## Download binary

https://github.com/wtetsu/gaze/releases

# How to use Gaze

The top priority of the Gaze's design is "easy to invoke".

By this command, Gaze starts watching the files in the current directory.

```
gaze .
```

On another terminal, run `vi a.py` and edit it. Gaze executes a.py in response to your file modifications!

### Other examples:

Gaze at one file.

```
gaze a.py
```

Gaze at subdirectories. Runs a modified file.

```
gaze 'src/**/*.rb'
```

Gaze at subdirectories. Runs a command to a modified file.

```
gaze 'src/**/*.js' -c "eslint {{file}}"
```

Kill an ongoing process, every time before it runs the next(Useful when you are writing servers)

```
gaze -r server.py
```

Kill an ongoing process, after 1000(ms)(Useful when you like to write infinite loops)

```
gaze -t 1000 complicated.py
```

### Configuration

Gaze is Language-agnostic.

But it has useful default configurations for several languages.

Due to the default configurations, the command below is valid.

```
gaze a.py
```

By default, it is the same as:

```
gaze a.py -c 'python "{{file}}"'
```

Gaze searches a configuration file according to it's priority rule.

1. A file specified by -f option
1. ./gaze.yml
1. ~/.gaze.yml
1. (Default)

You can display the default configuration by running `gaze -y`.

```yaml
commands:
  - ext: .go
    cmd: go run "{{file}}"
  - ext: .py
    cmd: python "{{file}}"
  - ext: .rb
    cmd: ruby "{{file}}"
  - ext: .js
    cmd: node "{{file}}"
  - ext: .d
    cmd: dmd -run "{{file}}"
  - ext: .groovy
    cmd: groovy "{{file}}"
  - ext: .php
    cmd: php "{{file}}"
  - ext: .pl
    cmd: perl "{{file}}"
  - ext: .java
    cmd: java "{{file}}"
  - ext: .kts
    cmd: kotlinc -script "{{file}}"
  - re: ^Dockerfile$
    cmd: docker build -f "{{file}}" .
```

You're able to have your own configuration very easily.

```
gaze -y > ~/.gaze.yml
vi ~/.gaze.yml
```

### Options:

```
Usage: gaze [options...] file(s)

Options:
  -c  A command string.
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -f  Specify a YAML configuration file.
  -v  Verbose mode.
  -q  Quiet mode.
  -y  Output the default configuration
  -h  Display help
  --color    Color(0:plain, 1:colorful)
  --version  Output version information

Examples:
  gaze .
  gaze *.rb
  gaze main.go
  gaze -c make '**/*.c'
  gaze -c "eslint {{file}}" 'src/**/*.js'
  gaze -r server.py
  gaze -t 1000 complicated.py

```

### Command format

You can specify a mustache style template as a command.

```
gaze -c 'echo "{{file}}" "{{abs}}"' .
```

| Parameter | Example                 |
| --------- | ----------------------- |
| {{file}}  | src/mod1/main.py        |
| {{ext}}   | .py                     |
| {{base}}  | main.py                 |
| {{base0}} | main                    |
| {{dir}}   | src/mod1                |
| {{abs}}   | /my/source/mod1/main.py |

# Third-party data

https://www.iconfinder.com/icons/2303106/eye_opened_public_visible_watch_icon
