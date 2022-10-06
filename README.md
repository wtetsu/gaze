
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

It greatly helps you focus on writing code!
![gaze02](https://user-images.githubusercontent.com/515948/73607575-1fbfe900-45fb-11ea-813e-6be6bf9ece6d.gif)

---

The usage of Gaze is quite simple.

```
gaze .
```

And invoke your favorite editor on another terminal and edit it!

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

👁️Gaze runs a script, **right after** you save it (e.g. Python, Ruby),

You can also use Gaze for these purposes:

- 👁️Gaze runs tests, **right after** you save a Ruby script
- 👁️Gaze runs linter, **right after** you save a JavaScript file
- 👁️Gaze runs "docker build .", **right after** you save Dockerfile
- 👁️And so forth...

---

Software development often forces us to execute the same command again and again, by hand!

Let's say, you started writing a really really simple Python script. You created a.py, wrote 5 lines of code and run "python a.py".
Since the result was not perfect, you edited a.py again, and run "python a.py" again.

Again and again...

Then, you found yourself going back and forth between the editor and terminal and typing the same command thousands of times.

That's totally waste of time and energy!

---

👁️Gaze runs a command on behalf of you, **right after** you edit files.

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

I developed Gaze in order to deal with my every day's coding.

Even though there are already many "update-and-run" type of tools, I would say Gaze is the best tool for quick coding because all the technical design decisions have been made for that purpose.

# How to use Gaze

The top priority of the Gaze's design is "easy to invoke".

By this command, Gaze starts watching the files in the current directory.

```
gaze .
```

On another terminal, run `vi a.py` and edit it. Gaze executes a.py in response to your file modifications!

### Other examples

Gaze at one file. You can just simply specify file names.

```
gaze a.py
```

---

Gaze doesn't have special options to specify files. You can use wildcards (\*, \*\*, ?) that shell users are familiar with. **You don't have to remember Gaze-specific command-line options!**

```
gaze "*.py"
```

---

Gaze at subdirectories. Runs a modified file.

```
gaze "src/**/*.rb"
```

---

Gaze at subdirectories. Runs a command to a modified file.

```
gaze "src/**/*.js" -c "eslint {{file}}"
```

---

Kill an ongoing process, every time before it runs the next. This is useful when you are writing a server.

```
gaze -r server.py
```

---

Kill an ongoing process, after 1000(ms). This is useful if you love to write infinite loops.

```
gaze -t 1000 complicated.py
```

---

In order to run multiple commands for one update, just simply write multiple lines (use quotations for general shells). If an exit code was not 0, Gaze doesn't invoke the next command.

```
gaze "*.cpp" -c "gcc {{file}} -o a.out
ls -l a.out
./a.out"
```

Here is output when a.cpp was updated.

```
[gcc a.cpp -o a.out](1/3)

[ls -l a.out](2/3)
-rwxr-xr-x 1 user group 42155 Mar  3 00:31 a.out

[./a.out](3/3)
hello, world!
```

When compilation failed:

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

But it has useful default configurations for several major languages (e.g. Go, Python, Ruby, JavaScript, D, Groovy, PHP, Java, Kotlin, Rust, C++, TypeScript, and Docker).

Thanks to the default configurations, the command below is valid. You don't have to specify python command.

```
gaze a.py
```

By default, the above is the same as:

```
gaze a.py -c 'python "{{file}}"'
```

Gaze searches a configuration file according to it's priority rule.

1. A file specified by -f option
1. ~/.config/gaze/gaze.yml
1. ~/.gaze.yml
1. (Default)

You can display the default YAML configuration by running `gaze -y`.

```yaml
commands:
  - ext: .bash
    cmd: bash "{{file}}"
  - ext: .cpp
    cmd: |
      gcc "{{file}}" -o"{{base0}}.out"
      ./"{{base0}}.out"
  - ext: .d
    cmd: dmd -run "{{file}}"
  - ext: .go
    cmd: go run "{{file}}"
  - ext: .groovy
    cmd: groovy "{{file}}"
  - ext: .java
    cmd: java "{{file}}"
  - ext: .js
    cmd: node "{{file}}"
  - ext: .kts
    cmd: kotlinc -script "{{file}}"
  - ext: .php
    cmd: php "{{file}}"
  - ext: .py
    cmd: python "{{file}}"
  - ext: .rb
    cmd: ruby "{{file}}"
  - ext: .rs
    cmd: |
      rustc "{{file}}" -o"{{base0}}.out"
      ./"{{base0}}.out"
  - ext: .sh
    cmd: sh "{{file}}"
  - ext: .ts
    cmd: |
      tsc "{{file}}" --out "{{base0}}.out"
      node ./"{{base0}}.out"
  - re: ^Dockerfile$
    cmd: docker build -f "{{file}}" .
```

Note:

- To specify both ext and re for one cmd is prohibited
- cmd can have multiple commands. In [YAML](https://en.wikipedia.org/wiki/YAML#Basic_components), a **vertical line(|)** is used to express multiple lines

You're able to have your own configuration very easily.

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
