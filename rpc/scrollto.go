package rpc

type ScrollTo struct {
	ViewID string `json:"view_id"`
	Col    int    `json:"col"`
	Line   int    `json:"line"`
}
