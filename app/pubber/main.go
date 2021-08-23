package main

import (
	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	pubber := NewPubber(app)
	app.Register(pubber)

	app.Run()

	app.Wait()
}
