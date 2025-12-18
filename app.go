package main

import (
	"context"
	"fmt"

	"mic-toggle/internal/config"
	hotkey "mic-toggle/internal/hotkey"
	"mic-toggle/internal/mic"

	"github.com/energye/systray"
	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	config         *config.Config
	windowVisible  bool
	hotkeyListener *hotkey.HotkeyListener
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	a.config = cfg
	a.windowVisible = false // window is hidden initially

	// Start hotkey listener in background
	a.hotkeyListener = &hotkey.HotkeyListener{}
	a.hotkeyListener.Start(a.config, func() {
		fmt.Println("Hotkey pressed!")

		// Toggle microphone
		muted, err := mic.ToggleMic()
		if err != nil {
			fmt.Println("Failed to toggle mic:", err)
		} else {
			fmt.Println("Mic muted:", muted)
		}

		// Play beep if enabled
		if a.config.PlayBeep {
			if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
				fmt.Println("Beep failed:", err)
			}
		}

		// Show notification if enabled
		if a.config.ShowNotification {
			status := "unmuted"
			if muted {
				status = "muted"
			}
			if err := beeep.Notify("Mic Toggle", "Microphone "+status, ""); err != nil {
				fmt.Println("Notification failed:", err)
			}
		}
	})

	// Start system tray in background
	go a.RunTray()
}

// ShowWindow shows the Wails window
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowCenter(a.ctx)
	a.windowVisible = true
}

// HideWindow hides the Wails window
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	a.windowVisible = false
}

// IsWindowVisible returns true if the window is currently shown
func (a *App) IsWindowVisible() bool {
	return a.windowVisible
}

// RunTray initializes the system tray
func (a *App) RunTray() {
	systray.Run(func() {
		systray.SetTitle("Mic Toggle")
		systray.SetTooltip("Mic Toggle App")
		systray.SetIcon(icon)

		// Click event toggles window visibility
		systray.SetOnClick(func(menu systray.IMenu) {
			if a.IsWindowVisible() {
				a.HideWindow()
			} else {
				a.ShowWindow()
			}
		})

		// Right click menu
		systray.SetOnRClick(func(menu systray.IMenu) {
			menu.ShowMenu()
		})

		mShowHide := systray.AddMenuItem("Show/Hide", "Show or hide the window")
		mQuit := systray.AddMenuItem("Quit", "Quit the application")

		mShowHide.Click(func() {
			if a.IsWindowVisible() {
				a.HideWindow()
			} else {
				a.ShowWindow()
			}
		})

		mQuit.Click(func() {
			systray.Quit()
			runtime.Quit(a.ctx)
		})
	}, func() {
		// Cleanup on exit
		fmt.Println("Systray exited")
	})
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {}

// Greet returns a greeting
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetHotkey returns the hotkey from config
func (a *App) GetHotkey() string {
	return a.config.Hotkey
}

// SetHotkey updates the hotkey and saves config
func (a *App) SetHotkey(hotkey string) error {
	a.config.Hotkey = hotkey
	if err := config.Save(a.config); err != nil {
		return err
	}

	// Restart hotkey listener immediately
	if a.hotkeyListener != nil {
		a.hotkeyListener.Stop()
	}
	a.hotkeyListener.Start(a.config, func() {
		// Dynamic behavior: check user options
		if a.config.PlayBeep {
			beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		}
		if a.config.ShowNotification {
			beeep.Notify("Mic Toggle", "Hotkey pressed!", "")
		}
	})

	return nil
}

func (a *App) SetPlayBeep(enabled bool) error {
	a.config.PlayBeep = enabled
	return config.Save(a.config)
}

func (a *App) SetShowNotification(enabled bool) error {
	a.config.ShowNotification = enabled
	return config.Save(a.config)
}

func (a *App) GetPlayBeep() bool {
	return a.config.PlayBeep
}

func (a *App) GetShowNotification() bool {
	return a.config.ShowNotification
}
