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

type Hotkey struct {
	Ctrl  bool
	Shift bool
	Alt   bool
	Win   bool
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
		case "win", "windows":
			h.Win = true
		default:
			h.Key = KeyNameToVK(p)
		}
	}
	return h
}

func KeyNameToVK(name string) types.VKCode {
	name = strings.ToUpper(strings.TrimSpace(name))

	// Single letter keys A-Z
	if len(name) == 1 && name[0] >= 'A' && name[0] <= 'Z' {
		return types.VKCode(name[0])
	}

	// Number keys 0-9
	if len(name) == 1 && name[0] >= '0' && name[0] <= '9' {
		return types.VKCode(name[0])
	}

	// Function keys
	switch name {
	case "F1":
		return types.VK_F1
	case "F2":
		return types.VK_F2
	case "F3":
		return types.VK_F3
	case "F4":
		return types.VK_F4
	case "F5":
		return types.VK_F5
	case "F6":
		return types.VK_F6
	case "F7":
		return types.VK_F7
	case "F8":
		return types.VK_F8
	case "F9":
		return types.VK_F9
	case "F10":
		return types.VK_F10
	case "F11":
		return types.VK_F11
	case "F12":
		return types.VK_F12
	case "F13":
		return types.VK_F13
	case "F14":
		return types.VK_F14
	case "F15":
		return types.VK_F15
	case "F16":
		return types.VK_F16
	case "F17":
		return types.VK_F17
	case "F18":
		return types.VK_F18
	case "F19":
		return types.VK_F19
	case "F20":
		return types.VK_F20
	case "F21":
		return types.VK_F21
	case "F22":
		return types.VK_F22
	case "F23":
		return types.VK_F23
	case "F24":
		return types.VK_F24

	// Special keys
	case "SPACE", "SPACEBAR":
		return types.VK_SPACE
	case "ENTER", "RETURN":
		return types.VK_RETURN
	case "ESCAPE", "ESC":
		return types.VK_ESCAPE
	case "TAB":
		return types.VK_TAB
	case "BACKSPACE", "BACK":
		return types.VK_BACK
	case "DELETE", "DEL":
		return types.VK_DELETE
	case "INSERT", "INS":
		return types.VK_INSERT
	case "HOME":
		return types.VK_HOME
	case "END":
		return types.VK_END
	case "PAGEUP", "PGUP":
		return types.VK_PRIOR
	case "PAGEDOWN", "PGDN":
		return types.VK_NEXT
	case "UP":
		return types.VK_UP
	case "DOWN":
		return types.VK_DOWN
	case "LEFT":
		return types.VK_LEFT
	case "RIGHT":
		return types.VK_RIGHT
	case "CAPSLOCK", "CAPS":
		return types.VK_CAPITAL
	case "NUMLOCK":
		return types.VK_NUMLOCK
	case "SCROLLLOCK", "SCROLL":
		return types.VK_SCROLL
	case "PRINTSCREEN", "PRINT":
		return types.VK_SNAPSHOT
	case "PAUSE":
		return types.VK_PAUSE

	// Numpad keys
	case "NUMPAD0", "NUM0":
		return types.VK_NUMPAD0
	case "NUMPAD1", "NUM1":
		return types.VK_NUMPAD1
	case "NUMPAD2", "NUM2":
		return types.VK_NUMPAD2
	case "NUMPAD3", "NUM3":
		return types.VK_NUMPAD3
	case "NUMPAD4", "NUM4":
		return types.VK_NUMPAD4
	case "NUMPAD5", "NUM5":
		return types.VK_NUMPAD5
	case "NUMPAD6", "NUM6":
		return types.VK_NUMPAD6
	case "NUMPAD7", "NUM7":
		return types.VK_NUMPAD7
	case "NUMPAD8", "NUM8":
		return types.VK_NUMPAD8
	case "NUMPAD9", "NUM9":
		return types.VK_NUMPAD9
	case "MULTIPLY", "NUM*":
		return types.VK_MULTIPLY
	case "ADD", "NUM+":
		return types.VK_ADD
	case "SUBTRACT", "NUM-":
		return types.VK_SUBTRACT
	case "DECIMAL", "NUM.":
		return types.VK_DECIMAL
	case "DIVIDE", "NUM/":
		return types.VK_DIVIDE

	// Symbol keys
	case "SEMICOLON", ";":
		return types.VK_OEM_1
	case "SLASH", "/":
		return types.VK_OEM_2
	case "TILDE", "~", "`":
		return types.VK_OEM_3
	case "LEFTBRACKET", "[":
		return types.VK_OEM_4
	case "BACKSLASH", "\\":
		return types.VK_OEM_5
	case "RIGHTBRACKET", "]":
		return types.VK_OEM_6
	case "QUOTE", "'":
		return types.VK_OEM_7
	case "PLUS", "=":
		return types.VK_OEM_PLUS
	case "COMMA", ",":
		return types.VK_OEM_COMMA
	case "MINUS", "-":
		return types.VK_OEM_MINUS
	case "PERIOD", ".":
		return types.VK_OEM_PERIOD
	}

	return 0
}

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
		winKeys := map[types.VKCode]bool{types.VK_LWIN: true, types.VK_RWIN: true}

		modifiers := map[string]bool{
			"ctrl":  false,
			"shift": false,
			"alt":   false,
			"win":   false,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case k := <-keyboardChan:
				if ctrlKeys[k.VKCode] {
					modifiers["ctrl"] = k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN
				}
				if shiftKeys[k.VKCode] {
					modifiers["shift"] = k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN
				}
				if altKeys[k.VKCode] {
					modifiers["alt"] = k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN
				}
				if winKeys[k.VKCode] {
					modifiers["win"] = k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN
				}

				if (k.Message == types.WM_KEYDOWN || k.Message == types.WM_SYSKEYDOWN) &&
					modifiers["ctrl"] == hotkey.Ctrl &&
					modifiers["shift"] == hotkey.Shift &&
					modifiers["alt"] == hotkey.Alt &&
					modifiers["win"] == hotkey.Win &&
					k.VKCode == hotkey.Key {
					callback()
				}
			}
		}
	}()

	return nil
}

func (h *HotkeyListener) Stop() {
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil
	}
}
