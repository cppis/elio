# elio  
`elio` means event loop io library.  
This library allows you to quickly write epoll-based servers.  

[![github-license](https://img.shields.io/github/license/cppis/elio)](https://img.shields.io/github/license/cppis/elio)
[![Build App](https://github.com/cppis/elio/actions/workflows/build-app.yml/badge.svg)](https://github.com/cppis/elio/actions/workflows/build-app.yml/badge.svg)
[![Publish App](https://github.com/cppis/elio/actions/workflows/publish-app.yml/badge.svg?tag=v0.1.7)](https://github.com/cppis/elio/actions/workflows/publish-app.yml)
[![github-license](https://img.shields.io/github/go-mod/go-version/cppis/elio)](https://img.shields.io/github/go-mod/go-version/cppis/elio)
[![tag version](https://img.shields.io/github/v/tag/cppis/elio)](https://img.shields.io/github/v/tag/cppis/elio)

<br/>

## Installation  
### Download `elio`  
```shell
$ git clone https://github.com/cppis/elio
$ cd elio
```

> Now, **$PWD** is the root path.  

<br/>

### [Setting `Skaffold` on Windows](docs/setting.skaffold.md)  
`Skaffold` settings on windows for continuous developing a Kubernetes-native app.  

<br/>

### [Setting `Helm` on Windows](docs/setting.helm.md)  
`Helm` settings on windows for managing packages for a Kubernetes.  

<br/><br/><br/>

## Run Echo  
`elio` use the [`viper`](https://github.com/spf13/viper) package to read configs.  
It reads configs from environment variable or *yaml/json* files.  
There is a predefined format of configs.  
for example, `{Service}_IN_URL` is a listen url of service.  

> The `elio.Service` interface has a `Name()` method,  
> which returns a string that is used as a prefix for environment variables.

You can run app in 3 ways.  
using `go run`, `docker` and `Skaffold`.   

first, move to echo project path.
```shell
cd app/echo
```

<br/>

### using `go run`  
To run `echo` service, run the following command:  
```shell
$ ECHO_IN_URL="0.0.0.0:7000" go run main.go
```

You can change the url of service `echo` by changing  
environment variable `ECHO_IN_URL`.

<br/>

### using `docker`  
To build `echo` image, run the following command:  
```shell
$ docker build -t elio:v0.1.4 .
$ docker tag elio:latest elio:v0.1.4
```

To run `echo` image, run the following command:  
```shell
$ docker run -d -e ECHO_IN_URL="0.0.0.0:7000" -p 7000:7000 -p 2345:2345 elio:v0.1.4
```

<br/>

### using `Skaffold`  
To use the `Skaffold`, you need thd following the [Setup `Skaffold`](#setup-skaffold).  
First move to Project root directory.  
```shell
$ cd {Project Root}
```

To run `echo` using `Skaffold`, run the following command:  
```shell
$ skaffold -f app\echo\k8s-resources\skaffold.yaml dev -p dev
```

Or, to run `echo` in debugging mode using `Skaffold`, run the following command:  
```shell
$ skaffold -f app\echo\k8s-resources\skaffold.yaml debug -p debug
```

<br/><br/><br/>

## Echo example  
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

<br/><br/><br/>

## Test  
You can test echo easily by using telnet.  
And, you can end server by send `q` character.  
