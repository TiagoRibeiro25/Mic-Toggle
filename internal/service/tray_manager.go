package service

import (
	"context"
	"fmt"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TrayManager struct {
	ctx           context.Context
	windowManager *WindowManager
	icon          []byte
}

func NewTrayManager(ctx context.Context, windowManager *WindowManager, icon []byte) *TrayManager {
	return &TrayManager{
		ctx:           ctx,
		windowManager: windowManager,
		icon:          icon,
	}
}

func (t *TrayManager) Run() {
	systray.Run(t.onReady, t.onExit)
}

func (t *TrayManager) onReady() {
	systray.SetTitle("Mic Toggle")
	systray.SetTooltip("Mic Toggle App")
	systray.SetIcon(t.icon)

	systray.SetOnClick(func(menu systray.IMenu) {
		t.windowManager.Toggle()
	})

	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})

	t.setupMenu()
}

func (t *TrayManager) setupMenu() {
	mShowHide := systray.AddMenuItem("Show/Hide", "Show or hide the window")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	mShowHide.Click(func() {
		t.windowManager.Toggle()
	})

	mQuit.Click(func() {
		systray.Quit()
		runtime.Quit(t.ctx)
	})
}

func (t *TrayManager) onExit() {
	fmt.Println("Systray exited")
}
