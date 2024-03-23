package vm

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type VM struct{}

// NewVM creates a new instance of the VM.
func NewVM() *VM {
	return &VM{}
}

// Run executes the VM logic.
// It starts the computer, initializes the program counter (PC),
// and continuously fetches and executes instructions from memory.
func (c *VM) Run() {
	var PCStart uint16 = 0x3000

	log.Println("Computer starting...")

	Register[R_PC] = PCStart

	fmt.Println()

	for {
		instruction := memory.Read(Register[R_PC])
		opcode := instruction >> 12
		Register[R_PC]++

		callOpcode(opcode, instruction)
	}
}

// LoadProgramImage loads a program image from the specified file path into the VM's memory.
// It reads the contents of the file, interprets them as a series of 16-bit values, and stores them in memory.
// The first two bytes of the file are treated as the origin address, and subsequent 16-bit values are stored
// in memory starting from that address.
// If an error occurs while reading the file or storing the values in memory, an error is returned.
func (vm *VM) LoadProgramImage(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Can't load %s: %s", path, err)
		return err
	}

	origin := binary.BigEndian.Uint16(bytes[:2])
	for i := 2; i < len(bytes); i += 2 {
		memory[origin] = binary.BigEndian.Uint16(bytes[i : i+2])
		origin++
	}
	return nil
}
