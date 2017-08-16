# kod (WIP, not ready)

[![GoDoc](https://godoc.org/github.com/linde12/kod/cmd?status.svg)](https://godoc.org/github.com/linde12/kod/cmd)

kod aims to be a fast and modern terminal code-editor. It is inspired by both vim & micro.

# Goals:
- Performance, it should never freeze or be slow
- Modes for efficiency, similar to vim
- Plugins(in lua, alternatively something similar to what xi-editor does)
- Modular(separate backend and frontend), so that different frontends may be added
- Small core, should be plugin driven

# Non-goals:
- Maintain a huge platform support(like vim), it will be supported by the major platforms and architectures

This is currently more of a proof-of-concept thing to get a basic understanding of how things can be done. Later i might fork micro and work from that, depending on how things go.

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

# License
MIT
