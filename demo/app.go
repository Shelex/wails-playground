package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	notify "github.com/willdot/gomacosnotify"
)

// App struct
type App struct {
	ctx     context.Context
	notify  *notify.Notifier
	version string
}

// NewApp creates a new App application struct
func NewApp() *App {
	notifier, _ := notify.New()
	return &App{
		version: "1.2.0",
		notify:  notifier,
	}
}

func (a *App) domReady(ctx context.Context) {
	a.updateCheckDialog()
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetVersion() string {
	return a.version
}

func (a *App) Greet(name string) string {
	message := fmt.Sprintf("Hello %s, It's show time!", name)
	a.Notify("Greeting", message)
	return message
}

func (a *App) Notify(title string, message string) error {
	runtime.LogInfof(a.ctx, "start notify for %s %s", title, message)

	notification := notify.Notification{
		Title:     "Demo-UI",
		SubTitle:  title,
		Message:   message,
		CloseText: "meh",
	}

	_ = notification.SetTimeout(10)

	_, err := a.notify.Send(notification)
	if err != nil {
		return err
	}

	return nil
}
