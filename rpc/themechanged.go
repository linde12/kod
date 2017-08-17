package rpc

type RGBA struct {
	R int32 `json:"r"`
	G int32 `json:"g"`
	B int32 `json:"b"`
	A int32 `json:"a"`
}

func (rgba RGBA) ToRGBInt() int32 {
	rgb := rgba.R<<16 + rgba.G<<8 + rgba.B
	return rgb
}

func (rgba RGBA) ToRGB() (int32, int32, int32) {
	return rgba.R, rgba.G, rgba.B
}

type Theme struct {
	Accent            *RGBA `json:"accent"`
	ActiveGuide       *RGBA `json:"active_guide"`
	Bg                *RGBA `json:"background"`
	Fg                *RGBA `json:"foreground"`
	BracketContentsFg *RGBA `json:"bracket_contents_foreground"`
	BracketsFg        *RGBA `json:"brackets_foreground"`
	Caret             *RGBA `json:"caret"`
	// TODO: Add all and document
}

type ThemeChanged struct {
	Name  string `json:"name"`
	Theme Theme  `json:"theme"`
}
