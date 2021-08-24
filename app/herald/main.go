package main

import (
	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	herald := NewHerald(app)
	app.Register(herald)

	app.Run()

	app.Wait()
}
