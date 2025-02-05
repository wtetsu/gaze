
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


# 👁️Gaze: Save & Run


Focus on your code, not commands!

<img src="https://github.com/user-attachments/assets/a7be2e53-b516-419e-afea-735fa7bc095e" width="800px" />

----


Repetitive command execution after each edit is a common frustration that disrupts our development flow. 😵‍💫

Let Gaze handle it!

- Save a.py -> 👁️Runs `python a.py`
- Save a.rb -> 👁️Runs `rubocop`
- Save a.go -> 👁️Runs `make build`
- Save Dockerfile -> 👁️Runs `docker build`
- And so forth...

## Installation

### Homebrew (macOS)

```
brew install gaze
```

Or, [download binary](https://github.com/wtetsu/gaze/releases)

## Quick start


Setting up Gaze is easy.

```
gaze .
```

Then, open your favorite editor in another terminal and start editing!


```
vi a.py
```



## Why Gaze? (Features)

Gaze is designed as a CLI tool that accelerates your coding.

- 📦 Easy to use, out-of-the-box
- ⚡ Lightning-fast response
- 🌎 Language-agnostic, editor-agnostic
- 🔧 Flexible configuration
- 📝 Create-and-rename file actions handling
- 🔍 Advanced options for more control
  - `-r`: Restart mode (useful for server applications)
  - `-t 2000`: Timeout (useful for preventing infinite loops)
- 🚀 Optimal parallel handling
  - See also: [Parallel handling](/doc/parallel.md)
  - <img src="doc/img/p04.png" width="300">
- 💻 Multiplatform (macOS, Windows, Linux)

---

Gaze was developed for supporting daily coding.

While many "update-and-run" tools exist, Gaze stands out with its focus on accelerating your coding workflow through a carefully considered technical design.

# How to use

Gaze prioritizes ease of use with its simple invocation.

```
gaze .
```

Then, switch to another terminal and run `vi a.py`. Gaze executes a.py in response to your file modifications.

---

Gaze at one file.

```
gaze a.py
```

---

Specify files using pattern matching (\*, \*\*, ?, {, })

```
gaze "*.py"
```

```
gaze "src/**/*.rb"
```

```
gaze "{aaa,bbb}/*.{rb,py}"
```

---

Specify a custom command by `-c` option.

```
gaze "src/**/*.js" -c "eslint {{file}}"
```

---

Kill the previous process before launching a new process. This is useful if you are writing a server.

```
gaze -r server.py
```

---

Kill a running process after 1000(ms). This is useful if you love infinite loops.

```
gaze -t 1000 complicated.py
```

---

Specify multiple commands within quotes, separated by newlines.

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

Gaze will not execute subsequent commands if a command exits with a non-zero status.


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

Gaze is language-agnostic.

For convenience, it provides helpful default configurations for a variety of popular languages (e.g., Go, Python, Ruby, JavaScript, Rust, etc.).


```
gaze a.py
```

By default, this command is equivalent to `gaze a.py -c 'python "{{file}}"'` because the default configuration includes:

```yaml
commands:
- ext: .py
  cmd: python "{{file}}"
```



You can view the default YAML configuration using `gaze -y`.


<details>
<summary>The default configuration</summary>

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
- ext: .zig
  cmd: zig run "{{file}}"
- re: ^Dockerfile$
  cmd: docker build -f "{{file}}" .

log:
  start: "[{{{command}}}]{{step}}"
  end: "({{elapsed_ms}}ms)"
```

</details>

---

To customize your configuration, create your own configuration file:


```
gaze -y > ~/.gaze.yml
vi ~/.gaze.yml
```

Gaze searches for a configuration file in the following order:

1. A file specified by -f option
1. ~/.config/gaze/gaze.yml
1. ~/.gaze.yml
1. (Default)





### Options:

```
Usage: gaze [options...] file(s)

Options:
  -c  Command(s) to run when files are changed.
  -r  Restart mode. Sends SIGTERM to the ongoing process before invoking the next command.
  -t  Timeout(ms). Sends SIGTERM to the ongoing process after the specified time has elapsed.
  -f  Specify a YAML configuration file.
  -v  Verbose mode. Displays additional information.
  -q  Quiet mode. Suppresses normal output.
  -y  Displays the default YAML configuration.
  -h  Displays help.
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

For more information: https://github.com/wtetsu/gaze
```

### Command format

You can use [Mustache](<https://en.wikipedia.org/wiki/Mustache_(template_system)>) templates in your commands.

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
  - See [go.mod](https://github.com/wtetsu/gaze/blob/master/go.mod) and [license.zip](https://github.com/wtetsu/gaze/releases)
