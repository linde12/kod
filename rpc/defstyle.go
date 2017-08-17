package rpc

type RGBAInt uint64

// TODO: This is kind of related to themechanged... Oversee structure
func (rgba RGBAInt) ToRGBInt() int32 {
	alpha := int32(rgba >> 24 & 0xFF)
	r, g, b := rgba.ToRGB()

	n := (1-alpha)*r + g + b
	return int32(n)
}

func (rgba RGBAInt) ToRGB() (int32, int32, int32) {
	red := rgba >> 16 & 0xFF
	green := rgba >> 8 & 0xFF
	blue := rgba & 0xFF

	// TODO: Calculate alpha?
	return int32(red), int32(green), int32(blue)
}

type DefineStyle struct {
	ID      int     `json:"id"`
	FgColor RGBAInt `json:"fg_color"`
	BgColor RGBAInt `json:"bg_color"`
}
