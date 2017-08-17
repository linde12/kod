package editor

import (
	"log"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/rpc"
)

var defaultStyle tcell.Style

type StyleMap struct {
	sync.Mutex
	values map[int]tcell.Style
}

func (m *StyleMap) Get(key int) tcell.Style {
	m.Lock()
	defer m.Unlock()
	return m.values[key]
}

func (m *StyleMap) Set(key int, value tcell.Style) {
	m.Lock()
	defer m.Unlock()
	m.values[key] = value
}

func NewStyleMap() *StyleMap {
	// TODO: Move this somewhere else
	return &StyleMap{
		values: make(map[int]tcell.Style),
	}
}

func (m *StyleMap) DefineStyle(defstyle *rpc.DefineStyle) {
	var style tcell.Style

	// TODO Make rpc.DefineStyle a map so we can see if FgColor exists or not
	if defstyle.FgColor != 0 {
		r, g, b := defstyle.FgColor.ToRGB()
		fg := tcell.NewRGBColor(r, g, b)
		log.Printf("fg: %d, %d, %d", r, g, b)
		style = defaultStyle.Foreground(fg)
	}

	if defstyle.BgColor != 0 {
		bg := tcell.NewRGBColor(defstyle.BgColor.ToRGB())
		style = defaultStyle.Background(bg)
	}

	m.Set(defstyle.ID, style)
}
