package vm

import (
	"fmt"
	"log"
)

type CPU struct{}

// NewCPU creates a new instance of the CPU.
func NewCPU() *CPU {
	return &CPU{}
}

// Run executes the CPU logic.
// It starts the computer, initializes the program counter (PC),
// and continuously fetches and executes instructions from memory.
func (c *CPU) Run() {
	// Implement CPU logic here
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
