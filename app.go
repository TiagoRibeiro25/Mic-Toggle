package main

import (
	"context"
	"fmt"

	"mic-toggle/internal/config"

	internal_hotkey "mic-toggle/internal/hotkey"
)

// App struct
type App struct {
	ctx    context.Context
	config *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	a.config = cfg

	go func() {
		internal_hotkey.ListenHotkey(a.config, func() {
			fmt.Println("Hotkey pressed!")
			// TODO: show toast or toggle mic
		})
	}()
}

// go func() {
// 	internal_hotkey.DebugKeyboard()
// }()

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// No action needed yet
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Cleanup if needed
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetHotkey returns the current hotkey from config
func (a *App) GetHotkey() string {
	return a.config.Hotkey
}

// SetHotkey updates the hotkey and saves to disk
func (a *App) SetHotkey(hotkey string) error {
	a.config.Hotkey = hotkey
	return config.Save(a.config)
}
