package vm

import (
	"github.com/eiannone/keyboard"
	"log"
)

func keyboardRead() uint16 {
	symb, controlKey, err := keyboard.GetSingleKey()

	if controlKey == keyboard.KeyEsc || controlKey == keyboard.KeyCtrlC {
		log.Println("Pressed escaping")
		// handle escape
	}

	if err != nil {
		log.Printf("Error, %s", err)
	}
	return uint16(symb)
}
