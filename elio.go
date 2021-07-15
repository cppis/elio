package elio

import "sync"

var _app *App
var onceApp sync.Once

// Elio get elio
func Elio() *App {
	onceApp.Do(func() {
		_app = NewApp()
	})
	return _app
}
