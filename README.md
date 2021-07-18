# elio  
`elio` means event loop io library.  
This library allows you to quickly write epoll-based servers.  

[![github-license](https://img.shields.io/github/license/cppis/elio)](https://img.shields.io/github/license/cppis/elio)
[![build workflow](https://github.com/cppis/elio/actions/workflows/build-elio.yml/badge.svg)](https://github.com/cppis/elio/actions/workflows/build-elio.yml/badge.svg)
[![github-license](https://img.shields.io/github/go-mod/go-version/cppis/elio)](https://img.shields.io/github/go-mod/go-version/cppis/elio)
[![tag version](https://img.shields.io/github/v/tag/cppis/elio)](https://img.shields.io/github/v/tag/cppis/elio)

<br/>

## Installation  
```shell
$ git clone https://github.com/cppis/elio
$ cd elio
$ go mod vendor
```

<br/><br/><br/>

## Run  
To run `echo` service. Run this command:  
```shell
$ ECHO_IN_URL="0.0.0.0:7000" go run app/echo/main.go
```

You can change the url of service `echo` using environment variable `ECHO_IN_URL`.

<br/><br/><br/>

## Example  
`app/echo` is simple echo example using `elio`.  
Here is the *main* function:  
```go
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
```

<br/>

`Echo` is implementation of `elio` service interface:  
`elio` service interface has following methods:  

```golang
package elio

import (
	"context"
	"time"
)

// Service service
type Service interface {
	Name() string
	OnInit(ctx context.Context, cancel context.CancelFunc) error
	OnExit()
	OnOpen(s *Session) error
	OnClose(s *Session, err error)
	OnError(s *Session, err error)
	OnRead(s *Session, in []byte) int
	OnWrite(s *Session, out []byte)
	OnLoop(host *IoHost, t time.Time, d time.Duration)
}
```
If network event happens, `elio` calls proper event method of service.  

<br/>

`echo` event implementations:  
```golang
func (e *Echo) OnOpen(s *elio.Session) error {
	fmt.Printf("o")

	return nil
}

func (e *Echo) OnClose(s *elio.Session, err error) {
	fmt.Printf("c")
}

func (e *Echo) OnError(s *elio.Session, err error) {
	fmt.Printf("e")
}

func (e *Echo) OnRead(s *elio.Session, in []byte) (processed int) {
	fmt.Printf("+%d", len(in))

	s.Write(in)

	if 'q' == in[0] {
		elio.Elio().End()
	}

	return processed
}

func (e *Echo) OnWrite(s *elio.Session, out []byte) {
	fmt.Printf("-%d", len(out))
}
```

<br/>

Service also has loop callback for run logics:

```go
const (
	// defaultFetchLimit default fetch limit
	defaultFetchLimit int = 2000
)

func (e *Echo) OnLoop(host *elio.IoHost, t time.Time, d time.Duration) {
	//host.RunDivision(t, r.callbackDivision)

	_, _ = host.Dispatching(t, defaultFetchLimit)

	e.prev = t
}
```