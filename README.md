# What it is?

Are you a terminal loving dev/dev-ops/SRE/whatever who hates `cd`-ing all day through multiple projects? Are you constantly pissed that your company decided service orientation was a good thing? Fear not!

CDP(CD Project, no connection with the creators of The Witcher even though I'm a fan) is a small CLI program that helps you exactly with that. It traverses predefined directories set by you, looking for hints that a directory is a project (called `markers`, stuff like `node_modules`, `.git`, configurable by you from the config file and possibly from CLI flags in the future).

# Instructions

## How to install

1. Clone the repository.
2. Either:
   a. Run `make install`, which will build and move the executable to `/usr/local/bin/cdp` (requires sudo and a Unix-like OS).
   b. Use the golang toolchain to build the executable: `go build -o <whatever> cmd/main.go` and do something with it (Might work on Windows, not tested).

## Usage

```man
Usage:
  cdp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  edit        Open project in the editor.
  help        Help about any command
  last        Show the last project opened with cdp.
  mux         Open project in a mux session
  shell       Open a shell in the project's directory.

Flags:
      --config string   config file (default is $HOME/.config/cdp/config.yaml (default "c")
  -h, --help            help for cdp
  -l, --last            Change to the last project.

Use "cdp [command] --help" for more information about a command.
```
