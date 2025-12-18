package service

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WindowManager struct {
	ctx     context.Context
	visible bool
}

func NewWindowManager(ctx context.Context) *WindowManager {
	return &WindowManager{
		ctx:     ctx,
		visible: false,
	}
}

func (w *WindowManager) Show() {
	runtime.WindowShow(w.ctx)
	runtime.WindowCenter(w.ctx)
	w.visible = true
}

func (w *WindowManager) Hide() {
	runtime.WindowHide(w.ctx)
	w.visible = false
}

func (w *WindowManager) Toggle() {
	if w.visible {
		w.Hide()
	} else {
		w.Show()
	}
}

func (w *WindowManager) IsVisible() bool {
	return w.visible
}
