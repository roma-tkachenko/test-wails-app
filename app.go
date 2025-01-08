package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"test-app/backend/service"
)

// App struct
type App struct {
	ctx        context.Context
	authStatus bool
}

type User struct {
	name               string
	avatarUrl          string
	vipStatus          string
	experience         Experience
	additionAttributes AdditionUserAttributes
}

type Experience struct {
	lvl       string
	currentXp string
	maxXp     string
}

type AdditionUserAttributes struct {
	group         string
	totalXp       string
	positionInTop string
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
	a.authStatus = false
	runtime.EventsEmit(a.ctx, "authStatus", a.authStatus)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Authenticate returns a greeting for the given username
func (a *App) Authenticate(username string, password string) string {
	message := fmt.Sprintf("Hello %s, you enter password: %s.", username, password)
	fmt.Printf(message)
	a.authStatus = true
	runtime.EventsEmit(a.ctx, "authStatus", a.authStatus)
	return message
}

func (a *App) Login(username string, password string) string {
	// Перевірка обов'язкових полів
	if username == "" {
		return "Username is required"
	}
	if password == "" {
		return "Password is required"
	}

	// Виконання авторизації
	success, err := service.Authorize(a.ctx, username, password)

	// Відправка статусу авторизації на фронтенд
	runtime.EventsEmit(a.ctx, "authStatus", success)

	if err != nil {
		return "Authorization failed: " + err.Error()
	}

	if success {
		return "Authorization successful"
	}

	return "Authorization failed"
}

// ClaimReward returns a greeting for the given username
func (a *App) ClaimReward() string {
	message := fmt.Sprintf("Claim Reward Runed.")
	fmt.Printf(message)
	_, err := service.ClaimRewards(a.ctx)
	if err != nil {
		fmt.Printf("Помилка отримання ревордів: %v", err)
		return "Помилка отримання ревордів"
	}
	return message
}

// ClaimReward returns a greeting for the given username
func (a *App) BoostClub() string {
	message := fmt.Sprintf("Boost Club Runed.")
	fmt.Printf(message)
	service.StartProcessing(a.ctx)
	return message
}
