package main

import (
	"fmt"

	"github.com/cppis/elio"
)

func main() {
	app := elio.Elio()

	sample := NewSample(app)
	app.Register(sample)

	fmt.Println("begin elio")
	defer fmt.Println("end elio")

	fmt.Println("run elio")
	app.Run()

	// go func(app *elio.App) {
	// 	time.Sleep(3 * time.Second)
	// 	fmt.Printf("\nend app...\n")
	// 	app.End()
	// }(app)

	// go func(sample *Sample) {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Printf("\ncancel sample...\n")
	// 	sample.cancel()
	// }(sample)

	app.Wait()

	fmt.Println("exit...")
}
