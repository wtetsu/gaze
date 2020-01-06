**Gaze is gazing at you**

![GAZE](https://user-images.githubusercontent.com/515948/71816598-828a9700-30c6-11ea-92c8-ca0154e98794.png)

[![Build Status](https://travis-ci.com/wtetsu/gaze.svg?branch=master)](https://travis-ci.com/wtetsu/gaze)

# Gaze

## What is Gaze?

Gaze runs a command, **right after** you saved something.

## How is it useful?

Software development sometimes forces us to execute the same command again and again by hands.

Let's say, you started writing a really really really simple Python script. You created a.py, wrote 5 lines of code and run "python a.py".
Since the result was not perfect, you edited a.py again, and run "python a.py" again.

Again and again...

Before you realized, you've saved the same files and executes the same command thousands of times!

Gaze runs a command Instead of you, **right after** you edit files.

### Use cases:

🚀**Right after** you save a script(e.g. Python), Gaze runs the script

You can also use Gaze for these purposes:

- 🚀**Right after** you save a Ruby script, Gaze runs tests
- 🚀**Right after** you save JavaScript files, Gaze runs linter
- 🚀**Right after** you save Dockerfile, Gaze runs "docker build ."
- And so forth...

## How to use Gaze

The top priority of the Gaze's design is "easy to invoke".

You never forget the command-line options in order to use for the main use cases.
In order to invoke Gazem all you have to do to is type four characters and push Enter.

```
gaze
```

By default, gaze watch the current directory.

This command same as "gaze".

### Command-line option examples:

```
gaze a.py
# same as:
# gaze a.py -c "python"
```

```
gaze Dockerfile
# same as:
# gaze Dockerfile -c "docker build"
```

```
gaze .
```

### Options:

```
Usage: gaze [files...] [options...]

Options:
  -c  Command.
  -q  Quiet.
  -r  Recursive.
  -p  Parallel.
  -f  Filter.
```

#### Command

You can specify a mustache style template string as a command.

```

```

| Parameter |     |     |     |     |
| --------- | --- | --- | --- | --- |
| {{path}}  |     |     |     |     |
| bbb       |     |     |     |     |
| cccc      |     |     |     |     |

# Install

## Brew(Only OSX)

(TODO)

## Get executables

(TODO)

## Build from source code

(TODO)

# Third-party data

https://www.iconfinder.com/icons/2303106/eye_opened_public_visible_watch_icon
