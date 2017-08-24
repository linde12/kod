package app

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type Application struct {
	*views.Application
	screen tcell.Screen
}

func NewApplication(screen tcell.Screen) *Application {
	app := &Application{}
	app.Application = &views.Application{}
	app.screen = screen
	app.SetScreen(screen)

	return app
}

func (a *Application) Size() (int, int) {
	return a.screen.Size()
}
