package main

import (
	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	puber := NewPuber(app)
	app.Register(puber)

	app.Run()

	app.Wait()
}
