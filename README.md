# How to install
1. Clone the repository.
2. Either:
    a. Run `sudo make install`, which will build and move the executable to `/usr/local/bin/cdp` (requires sudo and a Unix-like OS).
    b. Use the golang toolchain to build the executable: `go build -o <whatever> cmd/main.go` and do something with it (Might work on Windows, not tested).
