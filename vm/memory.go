package vm

import (
	"math"
)

// vmMemory represents the memory of the LC-3 virtual machine.
// It is implemented as an array of 16-bit unsigned integers.
type vmMemory [math.MaxUint16 + 1]uint16

// Read reads the value stored at the given memory address.
// If the address is MR_KBSR (keyboard status register),
// it checks if a key is pressed and updates the register accordingly.
// It returns the value stored at the given address.
func (mem *vmMemory) Read(address uint16) uint16 {
	if address == MR_KBSR {
		checkKey := keyboardRead()

		if checkKey != 0 {
			mem[MR_KBSR] = 1 << 15
			mem[MR_KBDR] = checkKey
		} else {
			mem[MR_KBSR] = 0
		}
	}
	return mem[address]
}

// Write the given value to the specified memory address.
func (mem *vmMemory) Write(address, value uint16) {
	mem[address] = value
}

// memory represents the virtual machine's memory.
var memory vmMemory
