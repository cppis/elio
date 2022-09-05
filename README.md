# elio  
`elio` means event loop io library.  
This library allows you to quickly write epoll-based servers.  

[![github-license](https://img.shields.io/github/license/cppis/elio)](https://img.shields.io/github/license/cppis/elio)
[![Build App](https://github.com/cppis/elio/actions/workflows/build-app.yml/badge.svg)](https://github.com/cppis/elio/actions/workflows/build-app.yml/badge.svg)
[![Publish App](https://github.com/cppis/elio/actions/workflows/publish-app.yml/badge.svg?tag=v0.1.7)](https://github.com/cppis/elio/actions/workflows/publish-app.yml)
[![github-license](https://img.shields.io/github/go-mod/go-version/cppis/elio)](https://img.shields.io/github/go-mod/go-version/cppis/elio)
[![tag version](https://img.shields.io/github/v/tag/cppis/elio)](https://img.shields.io/github/v/tag/cppis/elio)

<br/><br/><br/>

## Apps written using `elio` library  
### 🚀 [Running `Echo`](app/echo/README.md)  
`Echo` is a simple echo server.  

<br/>

### 🚀 [Running `Herald`](app/herald/README.md)  
`Herald` is a simple MQTT pub/sub test client.  

* Docker Hub: [cppis/herald](https://hub.docker.com/repository/docker/cppis/herald)  

<br/><br/><br/>

## Under the Hood of `Echo`  

`app/echo` is simple echo server written by `elio`.  
This is the *main* function of `app/echo`:  
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

The `echo` generated by `NewEcho` is a implementation of `elio` service interface:  
`elio` service interface has following method signatures:  

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

If network I/O event happens, `elio` calls proper event method of service.  

<br/>

This is a event implementations of `echo` service:  
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

Service also has loop callback for running logics at regular intervals:  

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
