
<p align="center">
  <img width="500" src="https://user-images.githubusercontent.com/515948/179385932-48ea38a3-3bbb-4f45-8d68-63dc076e757d.png" alt="gaze logo" />
  <br/>
  Gaze is gazing at you
</p>

<p align="center">
  <a href="https://github.com/wtetsu/gaze/actions?query=workflow%3ATest"><img src="https://github.com/wtetsu/gaze/workflows/Test/badge.svg" alt="Test" /></a>
  <a href="https://goreportcard.com/report/github.com/wtetsu/gaze"><img src="https://goreportcard.com/badge/github.com/wtetsu/gaze" alt="Go Report Card" /></a>
  <a href="https://codeclimate.com/github/wtetsu/gaze/maintainability"><img src="https://api.codeclimate.com/v1/badges/bd322b9104f5fcd3e37e/maintainability" alt="Maintainability" /></a>
  <a href="https://codecov.io/gh/wtetsu/gaze"><img src="https://codecov.io/gh/wtetsu/gaze/branch/master/graph/badge.svg" alt="codecov" /></a>
  <a href="https://pkg.go.dev/github.com/wtetsu/gaze"><img src="https://pkg.go.dev/badge/github.com/wtetsu/gaze.svg" alt="Go Reference"></a>
</p>


# What is Gaze?

👁️Gaze runs a command, **right after** you save a file.

It greatly helps you to focus on writing code!
![gaze02](https://user-images.githubusercontent.com/515948/73607575-1fbfe900-45fb-11ea-813e-6be6bf9ece6d.gif)

---

Setting up Gaze is easy.

```
gaze .
```

Then, invoke your favorite editor on another terminal and edit it!

```
vi a.py
```

## Installation

### Brew (for macOS)

```
brew install gaze
```

Or, [download binary](https://github.com/wtetsu/gaze/releases)

## Use cases:

- 👁️Runs a script right after any .py file is modified
- 👁️Runs tests right after any .rb file is modified
- 👁️Runs linter right after any .js is modified
- 👁️Runs `docker build .` right after you save Dockerfile
- And so forth...

---

Software development often forces us to execute the same command again and again, by hand!

Let's say, you started writing a really really simple Python script. You created a.py, wrote 5 lines of code and run "python a.py".
Since the result was not perfect, you edited a.py again, and run "python a.py" again.

Again and again...

Then, you found yourself going back and forth between the editor and terminal and typing the same command thousands of times.

That's totally waste of time and energy🙄

---

👁️Gaze runs a command for you, **right after** you save a file.

## Why Gaze? (Features)

Gaze is designed as a CLI tool that accelerates your coding.

- Easy to use, out-of-the-box
- Super quick reaction
- Language-agnostic, editor-agnostic
- Flexible configuration
- Useful advanced options
  - `-r`: restart (useful for server applications)
  - `-t 2000`: timeout (useful if you sometimes write infinite loops)
- Multiplatform (macOS, Windows, Linux)
- Can deal with "create-and-rename" type of editor's save behavior
  - Super major editors like Vim and Visual Studio are such editors
- Appropriate parallel handling
  - See also: [Parallel handling](/doc/parallel.md)
  - <img src="doc/img/p04.png" width="300">

---

Gaze was developed in order to support daily coding.

Even though there are already many "update-and-run" type of tools, I would say Gaze is the best for quick coding because all the technical design decisions have been made for that purpose.

# How to use Gaze

The top priority of the Gaze's design is "easy to invoke".

```
gaze .
```

Then, switch to another terminal and run `vi a.py`. Gaze executes a.py in response to your file modifications.

### Other examples

Gaze at one file.

```
gaze a.py
```

---

 You can use wildcards (\*, \*\*, ?) to specify files.

```
gaze "*.py"
```

---

Files in deep subdirectories.

```
gaze "src/**/*.rb"
```

---

Specify an arbitrary command by `-c` option.

```
gaze "src/**/*.js" -c "eslint {{file}}"
```

---

Kill the previous one before launching a new process. This is useful if you are writing a server.

```
gaze -r server.py
```

---

Kill an ongoing process after 1000(ms). This is useful if you love infinite loops.

```
gaze -t 1000 complicated.py
```

---

Specify multiple commands by using quotations.

```
gaze "*.cpp" -c "gcc {{file}} -o a.out
ls -l a.out
./a.out"
```

Output when a.cpp was updated.

```
[gcc a.cpp -o a.out](1/3)

[ls -l a.out](2/3)
-rwxr-xr-x 1 user group 42155 Mar  3 00:31 a.out

[./a.out](3/3)
hello, world!
```

If a certain command exited with non-zero, Gaze doesn't invoke the next command.

```
[gcc a.cpp -o a.out](1/3)
a.cpp: In function 'int main()':
a.cpp:5:28: error: expected ';' before '}' token
   printf("hello, world!\n")
                            ^
                            ;
 }
 ~
exit status 1
```

### Configuration

Gaze is Language-agnostic.

For convenience, it has useful default configurations for some major languages (e.g. Go, Python, Ruby, JavaScript, Rust, and so forth)

Thanks to the default configurations, the command below is valid.

```
gaze a.py
```

The above command is equivalent to `gaze a.py -c 'python "{{file}}"'`.


You can display the default YAML configuration by `gaze -y`.

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
  - ext: .java
    cmd: java "{{file}}"
  - ext: .kts
    cmd: kotlinc -script "{{file}}"
  - ext: .rs
    cmd: |
      rustc "{{file}}" -o"{{base0}}.out"
      ./"{{base0}}.out"
  - ext: .cpp
    cmd: |
      gcc "{{file}}" -o"{{base0}}.out"
      ./"{{base0}}.out"
  - ext: .ts
    cmd: |
      tsc "{{file}}" --out "{{base0}}.out"
      node ./"{{base0}}.out"
  - re: ^Dockerfile$
    cmd: docker build -f "{{file}}" .
```

Note:

- To specify both ext and re for one cmd is prohibited
- cmd can have multiple commands. Use vertical line(|) to write multiple commands

Gaze searches a configuration file according to its priority rule.

1. A file specified by -f option
1. ~/.config/gaze/gaze.yml
1. ~/.gaze.yml
1. (Default)

If you want to customize it, please set up your own configuration file.

```
gaze -y > ~/.gaze.yml
vi ~/.gaze.yml
```

### Options:

```
Usage: gaze [options...] file(s)

Options:
  -c  Command(s).
  -r  Restart mode. Send SIGTERM to an ongoing process before invoking next.
  -t  Timeout(ms). Send SIGTERM to an ongoing process after this time.
  -f  Specify a YAML configuration file.
  -v  Verbose mode.
  -q  Quiet mode.
  -y  Display the default YAML configuration.
  -h  Display help.
  --color    Color mode (0:plain, 1:colorful).
  --version  Display version information.

Examples:
  gaze .
  gaze main.go
  gaze a.rb b.rb
  gaze -c make "**/*.c"
  gaze -c "eslint {{file}}" "src/**/*.js"
  gaze -r server.py
  gaze -t 1000 complicated.py

```

### Command format

You can write [Mustache](<https://en.wikipedia.org/wiki/Mustache_(template_system)>) templates for commands.

```
gaze -c "echo {{file}} {{ext}} {{abs}}" .
```

| Parameter | Example                   |
| --------- | ------------------------- |
| {{file}}  | src/mod1/main.py          |
| {{ext}}   | .py                       |
| {{base}}  | main.py                   |
| {{base0}} | main                      |
| {{dir}}   | src/mod1                  |
| {{abs}}   | /my/proj/src/mod1/main.py |

# Third-party data

- Great Go libraries
  - See [go.mod](https://github.com/wtetsu/gaze/blob/master/go.mod) and [license.json](https://github.com/wtetsu/gaze/actions/workflows/license.yml)
