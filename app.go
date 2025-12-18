package main

import (
	"context"
	"fmt"

	"mic-toggle/internal/mic"
	"mic-toggle/internal/service"
)

// App struct
type App struct {
	ctx           context.Context
	appService    *service.AppService
	windowManager *service.WindowManager
	trayManager   *service.TrayManager
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{
		appService: service.NewAppService(),
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err := a.appService.Initialize(ctx); err != nil {
		panic(err)
	}

	a.windowManager = service.NewWindowManager(ctx)
	a.trayManager = service.NewTrayManager(ctx, a.windowManager, icon)

	go a.trayManager.Run()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {}

// Window management methods
func (a *App) ShowWindow() {
	a.windowManager.Show()
}

func (a *App) HideWindow() {
	a.windowManager.Hide()
}

func (a *App) IsWindowVisible() bool {
	return a.windowManager.IsVisible()
}

// Configuration methods
func (a *App) GetHotkey() string {
	return a.appService.GetConfig().Hotkey
}

func (a *App) SetHotkey(hotkey string) error {
	return a.appService.UpdateHotkey(hotkey)
}

func (a *App) GetPlayBeep() bool {
	return a.appService.GetConfig().PlayBeep
}

func (a *App) SetPlayBeep(enabled bool) error {
	return a.appService.UpdatePlayBeep(enabled)
}

func (a *App) GetShowNotification() bool {
	return a.appService.GetConfig().ShowNotification
}

func (a *App) SetShowNotification(enabled bool) error {
	return a.appService.UpdateShowNotification(enabled)
}

// Microphone methods
func (a *App) GetMicState() (bool, error) {
	return mic.IsMuted()
}

// Legacy method (consider removing if not used)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
