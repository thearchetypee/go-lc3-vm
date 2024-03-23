package vm

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

const (
	TRAP_GETC  = 0x20 /* get character from keyboard, not echoed onto the terminal */
	TRAP_OUT   = 0x21 /* output a character */
	TRAP_PUTS  = 0x22 /* output a word string */
	TRAP_IN    = 0x23 /* get character from keyboard, echoed onto the terminal */
	TRAP_PUTSP = 0x24 /* output a byte string */
	TRAP_HALT  = 0x25 /* halt the program */
)

func callTrap(trapVector uint16) {
	switch trapVector & 0xff {
	case TRAP_GETC:
		trapGetc()
	case TRAP_OUT:
		trapOut()
	case TRAP_PUTS:
		trapPuts()
	case TRAP_IN:
		trapIn()
	case TRAP_PUTSP:
		trapPutsp()
	case TRAP_HALT:
		trapHalt()
	}
}

// trapGetc reads a character from the keyboard and stores it in R0.
func trapGetc() {
	Register[R_R0] = keyboardRead()
}

// trapOut prints the character stored in Register[R_R0] to the standard output.
func trapOut() {
	fmt.Printf("%c", Register[R_R0])
}

// trapPuts prints the null-terminated string starting at the memory address stored in R_R0.
func trapPuts() {
	address := Register[R_R0]
	for {
		character := memory.Read(address)
		if character == 0 {
			break
		}
		fmt.Printf("%c", character)
		address++
	}
}

// trapIn reads a character from the keyboard and stores it in Register[R_R0].
func trapIn() {
	fmt.Printf("Enter a character: ")
	input := keyboardRead()
	fmt.Printf("%c", input)
	Register[R_R0] = input
}

// trapPutsp prints the null-terminated string at the memory address stored in R_R0 register.
// It reads each word from memory, extracts the characters, and prints them until a null character is encountered.
func trapPutsp() {
	address := Register[R_R0]
	for {
		word := memory.Read(address)
		char1 := word & 0xFF
		if char1 == 0 {
			break
		}
		fmt.Printf("%c", char1)
		char2 := word >> 8
		if char2 == 0 {
			break
		}
		fmt.Printf("%c", char2)
		address++
	}
}

// trapHalt halts the virtual machine.
func trapHalt() {
	fmt.Println("--- HALTING VM ---")
	os.Exit(0)
}

// keyboardRead reads a single key from the keyboard and returns its ASCII value as a uint16.
// If the key is the escape key (ESC) or the Ctrl+C combination, the function will print a message and halt the program.
// If there is an error while reading the key, the function will print an error message.
func keyboardRead() uint16 {
	symb, controlKey, err := keyboard.GetSingleKey()

	if controlKey == keyboard.KeyEsc || controlKey == keyboard.KeyCtrlC {
		fmt.Println("Pressed escaping")
		trapHalt()
	}

	if err != nil {
		fmt.Printf("Error, %s", err)
	}
	return uint16(symb)
}
