package main

import (
	"fmt"

	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	echo := NewSample(app)
	app.Register(echo)

	fmt.Println("begin elio")
	defer fmt.Println("end elio")

	fmt.Println("run elio")
	app.Run()

	// go func(app *elio.App) {
	// 	time.Sleep(3 * time.Second)
	// 	fmt.Printf("\nend app...\n")
	// 	app.End()
	// }(app)

	// go func(echo *Echo) {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Printf("\ncancel echo...\n")
	// 	echo.cancel()
	// }(echo)

	app.Wait()

	fmt.Println("exit...")
}
