package hotkey

import (
	"context"
	"fmt"
	"log"
	"strings"

	"mic-toggle/internal/config"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// HotkeyListener can start/stop a hotkey listener
type HotkeyListener struct {
	cancel context.CancelFunc
}

// Hotkey struct represents a parsed hotkey combination
type Hotkey struct {
	Ctrl  bool
	Shift bool
	Alt   bool
	Key   types.VKCode
}

// ParseHotkey converts string like "Ctrl+Shift+M" into Hotkey
func ParseHotkey(combo string) Hotkey {
	h := Hotkey{}
	parts := strings.Split(combo, "+")
	for _, p := range parts {
		switch strings.TrimSpace(strings.ToLower(p)) {
		case "ctrl":
			h.Ctrl = true
		case "shift":
			h.Shift = true
		case "alt":
			h.Alt = true
		default:
			h.Key = KeyNameToVK(p)
		}
	}
	return h
}

func KeyNameToVK(name string) types.VKCode {
	name = strings.ToUpper(strings.TrimSpace(name))
	if len(name) == 1 && name[0] >= 'A' && name[0] <= 'Z' {
		return types.VKCode(name[0])
	}
	return 0
}

// Start begins listening for the hotkey
func (h *HotkeyListener) Start(cfg *config.Config, callback func()) error {
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel

	go func() {
		log.SetFlags(0)
		log.SetPrefix("hotkey: ")

		keyboardChan := make(chan types.KeyboardEvent, 100)
		if err := keyboard.Install(nil, keyboardChan); err != nil {
			log.Fatal(err)
		}
		defer keyboard.Uninstall()

		hotkey := ParseHotkey(cfg.Hotkey)
		fmt.Println("Global hotkey listener started:", cfg.Hotkey)

		ctrlKeys := map[types.VKCode]bool{types.VK_LCONTROL: true, types.VK_RCONTROL: true}
		shiftKeys := map[types.VKCode]bool{types.VK_LSHIFT: true, types.VK_RSHIFT: true}
		altKeys := map[types.VKCode]bool{types.VK_LMENU: true, types.VK_RMENU: true}

		modifiers := map[string]bool{
			"ctrl":  false,
			"shift": false,
			"alt":   false,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case k := <-keyboardChan:
				if ctrlKeys[k.VKCode] {
					modifiers["ctrl"] = k.Message == types.WM_KEYDOWN
				}
				if shiftKeys[k.VKCode] {
					modifiers["shift"] = k.Message == types.WM_KEYDOWN
				}
				if altKeys[k.VKCode] {
					modifiers["alt"] = k.Message == types.WM_KEYDOWN
				}

				if k.Message == types.WM_KEYDOWN &&
					modifiers["ctrl"] == hotkey.Ctrl &&
					modifiers["shift"] == hotkey.Shift &&
					modifiers["alt"] == hotkey.Alt &&
					k.VKCode == hotkey.Key {
					callback()
				}
			}
		}
	}()

	return nil
}

// Stop cancels the listener
func (h *HotkeyListener) Stop() {
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil
	}
}
