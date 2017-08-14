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
* Unmarshal known JSON requests into structs
* Implement update protocol
* Read `viewHeight` lines and fill buffer
* Refactor away unused code(e.g. lineArray, buffer_reader, cursor)
* Add basic editing functionality
* Add vertical scrolling
* Add horizontal scrolling(client-side only?)
* Tabs
* Views within tabs

# License
MIT
