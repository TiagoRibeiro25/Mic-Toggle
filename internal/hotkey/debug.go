//go:build windows
// +build windows

package hotkey

import (
	"fmt"
	"log"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
)

// DebugKeyboard logs every key pressed/released
func DebugKeyboard() {
	log.SetFlags(0)
	log.SetPrefix("debug: ")

	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Uninstall()

	fmt.Println("Keyboard debug listener started. Press any key...")

	for k := range keyboardChan {
		var action string
		switch k.Message {
		case types.WM_KEYDOWN:
			action = "KeyDown"
		case types.WM_KEYUP:
			action = "KeyUp"
		case types.WM_SYSKEYDOWN:
			action = "SysKeyDown"
		case types.WM_SYSKEYUP:
			action = "SysKeyUp"
		default:
			action = fmt.Sprintf("Unknown(%d)", k.Message)
		}

		fmt.Printf("VKCode: %d (%s) - %s\n", k.VKCode, types.VKCode(k.VKCode).String(), action)
	}
}
