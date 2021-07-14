package main

import (
	"fmt"
	"time"

	"github.com/cppis/elio"
)

func main() {
  e := elio.Elio()

  fmt.Println("begin elio")
  defer fmt.Println("end elio")

  fmt.Println("run elio")
  e.Run()

  go func() {
    time.Sleep(3*time.Second)
    e.End()
  }()

  e.Wait()
}
