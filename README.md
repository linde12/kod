# kod (WIP)

[![GoDoc](https://godoc.org/github.com/linde12/kod/cmd?status.svg)](https://godoc.org/github.com/linde12/kod)

![kod screenshot](/.github/scrot.png?raw=true)

kod aims to be a fast and modern terminal code-editor. It is inspired by both vim & micro. It's built using xi-editor as backend.

# Goals:
- Performance, it should *never* freeze or be slow
- Modes for efficient editing, similar to vim
- Plugins(via xi-core, an example is xi-syntect-plugin which is shown in the picture above)
- Modular, uses xi-editor as backend

# Non-goals:
- Maintain a huge platform support(like vim), it will be supported by the major platforms and architectures

# Installation
This assumes that you've built xi-core and your plugins(e.g. xi-syntect-plugin for syntax highlighting) for your platform and placed them in the same folder as the `kod` executable. Currently there has been no effort to make the editor run outside it's current working directory, so you will have to run it from the executables location. Feel free to add a PR to change this, otherwise i'll add it later on.

To compile `kod` you'll have to run `go build` from the root of the project. This will generate an executable for your system & arch. To build and test it on other platforms you can cross-compile by setting the environment variables `GOOS` and `GOARCH`(for more information about this i'd suggest you take a look at [this article](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5).

# TODO
- [ ] Unmarshal known JSON requests into structs, partially done
- [x] Implement update protocol
- [x] Read `viewHeight` lines and fill buffer
- [x] Refactor away unused code(e.g. lineArray, buffer_reader, cursor)
- [ ] Restructure and refactor when editor is more complete
- [x] Add basic editing functionality
- [x] Add vertical scrolling
- [ ] Add horizontal scrolling
- [x] Indentation (local, xi doesn't fully support yet AFAIK)
- [ ] Tabs
- [ ] Views within tabs
- [ ] Fix all TODOs in the code
- [ ] Respect alpha value in ARGB
- [ ] Fix mutexes on shared memory(Go 1.9 has concurrent maps builtin)
- [ ] Cleanup view and inputhandler
- [ ] Find(search) in file
- [ ] Support multiple cursors
- [ ] Make editor runnable outside current CWD
- [ ] Add tests
- [ ] Go plugin
- [ ] ??? I probably missed something, feel free to add a PR
- [ ] Make the editor stable enough to use as a daily driver
- [ ] Some kind of modal mode
- [ ] Display line numbers/gutter

# License
MIT
