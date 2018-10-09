# kod (WIP)

[![GoDoc](https://godoc.org/github.com/linde12/kod/cmd?status.svg)](https://godoc.org/github.com/linde12/kod)

![kod screenshot](/.github/screenshot.png?raw=true)

kod aims to be a fast and modern terminal code-editor. It is inspired by both vim & micro. It's built using xi-editor as backend.

# Goals:
- Performance, it should *never* freeze or be slow
- Modes for efficient editing, similar to vim
- Plugins (via xi-core, an example is xi-syntect-plugin which is shown in the picture above)
- Modular, uses xi-editor as backend

# Non-goals:
- Maintain a huge platform support (like vim), it will be supported by the major platforms and architectures

# Installation
kod expects `xi-core` to be set in your `$PATH`. For more information how to build xi-core, please have a look at [xi-editor's github page](https://github.com/google/xi-editor).

Once xi-core is available you should be able to install with following
```
 go get -u -v github.com/linde12/kod
```
To ensure that kod is available, run `which kod`. It prints kod path in GO binary.
If you see `kod not found`, you need to add your `$GOPATH/bin` to PATH
# TODO
- [ ] Unmarshal known JSON requests into structs, partially done
- [x] Implement update method
- [x] Read `viewHeight` lines and fill buffer
- [x] Refactor away unused code(e.g. lineArray, buffer_reader, cursor)
- [x] Add (very) basic editing functionality
- [x] Add vertical scrolling
- [ ] Add horizontal scrolling
- [x] Indentation (local, xi doesn't fully support yet AFAIK)
- [ ] Respect alpha value in ARGB
- [ ] Cleanup view and inputhandler
- [ ] Find(search) in file
- [ ] Support multiple cursors
- [x] Make editor runnable outside current CWD
- [x] Display line numbers/gutter (very basic)
- [ ] A lot of other things...

# License
MIT
