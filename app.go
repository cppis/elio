package elio

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// App app
type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config
	wg     sync.WaitGroup
}

// GetContext get context
func (a *App) GetContext() context.Context {
	return a.ctx
}

// End call cancel callback
func (a *App) End() {
	a.cancel()
}

// Wait wait app
func (a *App) Wait() {
	a.wg.Wait()
}

func (a *App) init() {
	a.config = NewConfig()
	a.ctx, a.cancel = context.WithCancel(context.Background())
}

// Run run
func (a *App) Run() {
	a.wg.Add(1)

	go func(i time.Duration) {
		defer func() {
			a.wg.Done()
		}()
	
		ticker := time.NewTicker(i)
		for {
			select {
			case <-ticker.C:
				// TODO: work in here
				fmt.Printf(".")
			case <-a.ctx.Done():
				return
			}
		}
	}(100 * time.Millisecond)
}

func NewApp() *App {
	a := new(App)
	if nil != a {
		a.init()
	}
	return a
}

var _app *App
var onceApp sync.Once

// Elio get elio
func Elio() *App {
	onceApp.Do(func() {
		_app = NewApp()
	})
	return _app
}
