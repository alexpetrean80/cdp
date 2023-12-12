# What it is?

Are you a terminal loving dev/dev-ops/SRE/whatever who hates `cd`-ing all day through multiple projects? Are you constantly pissed that your company decided service orientation was a good thing? Fear not!

CDP(CD Project, no connection with the creators of The Witcher even though I'm a fan) is a small CLI program that helps you exactly with that. It traverses predefined directories set by you, looking for hints that a directory is a project (called `markers`, stuff like `node_modules`, `.git`, configurable by you from the config file and possibly from CLI flags in the future).

# Instructions

## How to install

1. Clone the repository.
2. Either:
   a. Run `sudo make install`, which will build and move the executable to `/usr/local/bin/cdp` (requires sudo and a Unix-like OS).
   b. Use the golang toolchain to build the executable: `go build -o <whatever> cmd/main.go` and do something with it (Might work on Windows, not tested).

## Usage

```man
NAME:
   cdp - Move between projects seamlessly

USAGE:
   cdp [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --edit, -e     Open the project in the configured editor. (mutually exclusive with -t) (default: false)
   --browser, -b  Open the project in the browser. github-cli required (default: false)
   --latest, -l   Open the latest project (default: false)
   --tmux, -t     Open the project in a new tmux session. (mutually exclusive with -o) tmux required. (default: false)
   --help, -h     show help
```
