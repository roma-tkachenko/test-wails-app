package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	authStatus bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	fmt.Printf("NewApp\n")
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.authStatus = false
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Authenticate returns a greeting for the given username
func (a *App) Authenticate(username string, password string) string {
	message := fmt.Sprintf("Hello %s, you enter password: %s.", username, password)
	fmt.Printf(message)
	runtime.EventsEmit(a.ctx, "authStatus", a.authStatus)
	return message
}
