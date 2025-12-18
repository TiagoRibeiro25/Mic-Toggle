package main

import (
	"embed"
	"fmt"
	"os"
	"syscall"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"golang.org/x/sys/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Named mutex to allow only one instance
	mutexName := "MicToggleSingletonMutex"
	mutex, err := windows.CreateMutex(nil, false, syscall.StringToUTF16Ptr(mutexName))
	if err != nil {
		fmt.Println("Error creating mutex:", err)
		os.Exit(1)
	}
	lastErr := windows.GetLastError()
	if lastErr == windows.ERROR_ALREADY_EXISTS {
		fmt.Println("Another instance is already running. Exiting...")
		os.Exit(0)
	}
	defer windows.ReleaseMutex(mutex)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:            "Mic Toggle",
		Width:            340,
		Height:           500,
		DisableResize:    true,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		StartHidden:       true,
		HideWindowOnClose: true,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
