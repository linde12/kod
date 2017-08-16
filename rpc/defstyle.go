package rpc

type RGBA int64

func (rgba RGBA) ToRGB() int32 {
	alpha := rgba << 24
	rgb := rgba >> 8
	return int32((1 - alpha) * rgb)
}

type DefineStyle struct {
	ID      int  `json:"id"`
	FgColor RGBA `json:"fg_color"`
	BgColor RGBA `json:"bg_color"`
}
