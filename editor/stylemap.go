package editor

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/rpc"
)

var defaultStyle tcell.Style

type StyleMap map[int]tcell.Style

func NewStyleMap() StyleMap {
	defaultStyle = defaultStyle.Background(tcell.ColorBrown)
	return make(StyleMap)
}

func (sm StyleMap) DefineStyle(defstyle *rpc.DefineStyle) {
	var style tcell.Style
	log.Printf("color 0x%x", defstyle.FgColor.ToRGB())

	// TODO Make rpc.DefineStyle a map so we can see if FgColor exists or not
	if defstyle.FgColor != 0 {
		style = defaultStyle.Foreground(tcell.Color(defstyle.FgColor.ToRGB()))
	}

	if defstyle.BgColor != 0 {
		style = defaultStyle.Background(tcell.Color(defstyle.BgColor.ToRGB()))
	}

	sm[defstyle.ID] = style
}
