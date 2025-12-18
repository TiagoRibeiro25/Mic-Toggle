//go:build windows
// +build windows

package hotkey

import (
	"fmt"
	"log"
	"strings"

	"mic-toggle/internal/config"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// Hotkey struct represents a parsed hotkey combination
type Hotkey struct {
	Ctrl  bool
	Shift bool
	Alt   bool
	Key   types.VKCode
}

// ParseHotkey converts a string like "Ctrl+Shift+M" into a Hotkey struct
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

// KeyNameToVK converts a key name (e.g., "M") to a VKCode
func KeyNameToVK(name string) types.VKCode {
	name = strings.ToUpper(strings.TrimSpace(name))
	if len(name) == 1 && name[0] >= 'A' && name[0] <= 'Z' {
		return types.VKCode(name[0])
	}
	// TODO: Extend for numbers, function keys, etc.
	return 0
}

// ListenHotkey installs a global hook and calls callback when the hotkey is pressed
func ListenHotkey(cfg *config.Config, callback func()) {
	log.SetFlags(0)
	log.SetPrefix("hotkey: ")

	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Uninstall()

	hotkey := ParseHotkey(cfg.Hotkey)
	fmt.Println("Global hotkey listener started:", cfg.Hotkey)

	// Map all left/right variants for modifier keys
	ctrlKeys := map[types.VKCode]bool{types.VK_LCONTROL: true, types.VK_RCONTROL: true}
	shiftKeys := map[types.VKCode]bool{types.VK_LSHIFT: true, types.VK_RSHIFT: true}
	altKeys := map[types.VKCode]bool{types.VK_LMENU: true, types.VK_RMENU: true}

	// Track modifier states
	modifiers := map[string]bool{
		"ctrl":  false,
		"shift": false,
		"alt":   false,
	}

	for k := range keyboardChan {
		// Update modifiers
		if ctrlKeys[k.VKCode] {
			modifiers["ctrl"] = k.Message == types.WM_KEYDOWN
		}
		if shiftKeys[k.VKCode] {
			modifiers["shift"] = k.Message == types.WM_KEYDOWN
		}
		if altKeys[k.VKCode] {
			modifiers["alt"] = k.Message == types.WM_KEYDOWN
		}

		// Check main key
		if k.Message == types.WM_KEYDOWN &&
			modifiers["ctrl"] == hotkey.Ctrl &&
			modifiers["shift"] == hotkey.Shift &&
			modifiers["alt"] == hotkey.Alt &&
			k.VKCode == hotkey.Key {

			callback()
		}
	}
}
