package service

import (
	"context"
	"fmt"
	"mic-toggle/internal/config"
	"mic-toggle/internal/hotkey"
	"mic-toggle/internal/mic"

	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppService struct {
	ctx            context.Context
	config         *config.Config
	hotkeyListener *hotkey.HotkeyListener
}

func NewAppService() *AppService {
	return &AppService{}
}

func (s *AppService) Initialize(ctx context.Context) error {
	s.ctx = ctx

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	s.config = cfg

	s.startHotkeyListener()
	return nil
}

func (s *AppService) startHotkeyListener() {
	s.hotkeyListener = &hotkey.HotkeyListener{}
	s.hotkeyListener.Start(s.config, func() {
		s.handleHotkeyPress()
	})
}

func (s *AppService) handleHotkeyPress() {
	fmt.Println("Hotkey pressed!")

	muted, err := mic.ToggleMic()
	if err != nil {
		fmt.Println("Failed to toggle mic:", err)
		return
	}

	fmt.Println("Mic muted:", muted)
	runtime.EventsEmit(s.ctx, "micStateChanged", muted)

	s.playFeedback(muted)
}

func (s *AppService) playFeedback(muted bool) {
	if s.config.PlayBeep {
		if err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration); err != nil {
			fmt.Println("Beep failed:", err)
		}
	}

	if s.config.ShowNotification {
		status := "unmuted"
		if muted {
			status = "muted"
		}
		if err := beeep.Notify("Mic Toggle", "Microphone "+status, ""); err != nil {
			fmt.Println("Notification failed:", err)
		}
	}
}

func (s *AppService) UpdateHotkey(hotkey string) error {
	s.config.Hotkey = hotkey
	if err := config.Save(s.config); err != nil {
		return err
	}

	if s.hotkeyListener != nil {
		s.hotkeyListener.Stop()
	}
	s.startHotkeyListener()

	return nil
}

func (s *AppService) UpdatePlayBeep(enabled bool) error {
	s.config.PlayBeep = enabled
	return config.Save(s.config)
}

func (s *AppService) UpdateShowNotification(enabled bool) error {
	s.config.ShowNotification = enabled
	return config.Save(s.config)
}

func (s *AppService) GetConfig() *config.Config {
	return s.config
}
