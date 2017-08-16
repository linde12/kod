package editor

import "github.com/linde12/kod/rpc"

type InputHandler struct {
	ViewID string
	// TODO: Maybe just keep this stuff on View struct
	FilePath string
	c        *rpc.Connection
}

func (ih *InputHandler) edit(params rpc.Object) {
	params["view_id"] = ih.ViewID
	if _, ok := params["params"]; !ok {
		params["params"] = &rpc.Object{}
	}

	ih.c.Notify(&rpc.Request{
		Method: "edit",
		ViewID: ih.ViewID,
		Params: params,
	})
}

func (ih *InputHandler) MoveLeft() {
	ih.edit(rpc.Object{"method": "move_left"})
}

func (ih *InputHandler) MoveRight() {
	ih.edit(rpc.Object{"method": "move_right"})
}

func (ih *InputHandler) MoveUp() {
	ih.edit(rpc.Object{"method": "move_up"})
}

func (ih *InputHandler) MoveDown() {
	ih.edit(rpc.Object{"method": "move_down"})
}

func (ih *InputHandler) DeleteBackward() {
	ih.edit(rpc.Object{"method": "delete_backward"})
}

func (ih *InputHandler) Tab() {
	ih.edit(rpc.Object{"method": "insert_tab"})
}

func (ih *InputHandler) Newline() {
	ih.edit(rpc.Object{"method": "insert_newline"})
}

func (ih *InputHandler) Insert(char string) {
	ih.c.Notify(&rpc.Request{
		Method: "edit",
		ViewID: ih.ViewID,
		Params: rpc.Object{"method": "insert", "params": rpc.Object{"chars": char}, "view_id": ih.ViewID},
	})
}

func (ih *InputHandler) Save() {
	ih.c.Notify(&rpc.Request{
		Method: "save",
		Params: rpc.Object{
			"view_id":   ih.ViewID,
			"file_path": ih.FilePath,
		},
	})
}
