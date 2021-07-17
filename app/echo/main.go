package main

import (
	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	echo := NewEcho(app)
	app.Register(echo)

	app.Run()

	app.Wait()
}
