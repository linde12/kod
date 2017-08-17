# kod (WIP)

[![GoDoc](https://godoc.org/github.com/linde12/kod/cmd?status.svg)](https://godoc.org/github.com/linde12/kod)

kod aims to be a fast and modern terminal code-editor. It is inspired by both vim & micro. It's built using xi-editor as backend.

# Goals:
- Performance, it should *never* freeze or be slow
- Modes for efficient editing, similar to vim
- Plugins(via xi-core)
- Modular, uses xi-editor as backend

# Non-goals:
- Maintain a huge platform support(like vim), it will be supported by the major platforms and architectures

# TODO
- [ ] Unmarshal known JSON requests into structs, partially done
- [x] Implement update protocol
- [x] Read `viewHeight` lines and fill buffer
- [x] Refactor away unused code(e.g. lineArray, buffer_reader, cursor)
- [ ] Restructure and refactor when editor is more complete
- [x] Add basic editing functionality
- [ ] Add vertical scrolling
- [x] Add horizontal scrolling
- [x] Indentation (local, xi doesn't fully support yet AFAIK)
- [ ] Tabs
- [ ] Views within tabs
- [ ] Fix all TODOs in the code
- [ ] Go plugin
- [ ] Make the editor stable enough to use as a daily driver
- [ ] ??? I probably missed something, feel free to add a PR

# License
MIT
