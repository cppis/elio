package main

import (
	"os"

	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	echo := NewEcho(app)
	app.Register(echo)

	os.Setenv("ECHO_IN_URL", "0.0.0.0:7000")

	app.Run()

	app.Wait()
}
