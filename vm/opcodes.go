package vm

import "log"

const (
	OP_BR   = iota /* branch */
	OP_ADD         /* add  */
	OP_LD          /* load */
	OP_ST          /* store */
	OP_JSR         /* jump register */
	OP_AND         /* bitwise and */
	OP_LDR         /* load register */
	OP_STR         /* store register */
	OP_RTI         /* unused */
	OP_NOT         /* bitwise not */
	OP_LDI         /* load indirect */
	OP_STI         /* store indirect */
	OP_JMP         /* jump */
	OP_RES         /* reserved (unused) */
	OP_LEA         /* load effective address */
	OP_TRAP        /* execute trap */
)

var OPCODES = map[uint16]func(uint16){
	OP_BR:   opBR,
	OP_ADD:  opADD,
	OP_LD:   opLD,
	OP_ST:   opST,
	OP_JSR:  opJSR,
	OP_AND:  opAND,
	OP_LDR:  opLDR,
	OP_STR:  opSTR,
	OP_RTI:  opRTI,
	OP_NOT:  opNOT,
	OP_LDI:  opLDI,
	OP_STI:  opSTI,
	OP_JMP:  opJMP,
	OP_RES:  opRES,
	OP_LEA:  opLEA,
	OP_TRAP: opTRAP,
}

// getR0 extracts the value of register R0 from the given instruction.
// It shifts the instruction 9 bits to the right and applies a bitwise AND operation
// with 0x7 to extract the desired bits.
func getR0(instruction uint16) uint16 {
	return (instruction >> 9) & 0x7
}

// getR1 extracts the value of register R1 from the given instruction.
// It shifts the instruction 6 bits to the right and applies a bitwise AND operation
// with 0x7 to extract the desired bits.
func getR1(instruction uint16) uint16 {
	return (instruction >> 6) & 0x7
}

func getImmFlag(instruction uint16) uint16 {
	return (instruction >> 5) & 0x1
}

// signExtend extends the sign of a given value to the specified number of bits.
// It takes an unsigned 16-bit value `x` and the number of bits `bitCount` to extend to.
// If the most significant bit of `x` is set (1), it performs sign extension by filling the additional bits with 1s.
// Otherwise, it leaves the additional bits as 0s.
// The resulting extended value is returned.
func signExtend(x uint16, bitCount int) uint16 {
	if ((x >> (bitCount - 1)) & 1) != 0 {
		x |= (0xFFFF << bitCount)
	}
	return x
}

func updateFlags(regIndex uint16) {
	if Register[regIndex] == 0 {
		Register[R_COND] = FL_ZRO
	} else if Register[regIndex]>>15 == 1 {
		Register[R_COND] = FL_NEG
	} else {
		Register[R_COND] = FL_POS
	}
}

// opBR performs a conditional branch based on the condition flag in the instruction.
// It updates the program counter (PC) if the condition flag matches the value in the register.
func opBR(instruction uint16) {
	var pcOffset = signExtend(instruction&0x1ff, 9)
	var condFlag = (instruction >> 9) & 0x7
	if condFlag&Register[R_COND] != 0 {
		Register[R_PC] += pcOffset
	}
}

// opADD performs the ADD operation.
// It takes an instruction as input and updates the appropriate registers based on the instruction.
// If the immediate flag is set, it adds the sign-extended immediate value to the value in register r1 and stores the result in register r0.
// If the immediate flag is not set, it adds the value in register r2 to the value in register r1 and stores the result in register r0.
// The function also updates the flags based on the result in register r0.
func opADD(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)

	if getImmFlag(instruction) == 1 {
		var imm5 = signExtend(instruction&0x1f, 5)
		Register[r0] = Register[r1] + imm5
	} else {
		var r2 = instruction & 0x7
		Register[r0] = Register[r1] + Register[r2]
	}

	updateFlags(r0)
}

// opLD loads a value from memory into a register.
// It takes an instruction as input and extracts the register number (r0) and the PC offset.
// The value at the memory address (PC + offset) is then loaded into the specified register (r0).
// Finally, the flags are updated based on the loaded value.
func opLD(instruction uint16) {
	var r0 = getR0(instruction)
	var pcOffset = signExtend(instruction&0x1ff, 9)
	Register[r0] = memory.Read(Register[R_PC] + pcOffset)
	updateFlags(r0)
}

// opST stores the value in the register r0 to memory at the address specified by the PC offset.
// It takes an instruction as input, extracts the r0 register and the PC offset, and writes the value in r0 to memory.
func opST(instruction uint16) {
	var r0 = getR0(instruction)
	var pcOffset = signExtend(instruction&0x1ff, 9)
	memory.Write(Register[R_PC]+pcOffset, Register[r0])
}

// opJSR performs the JSR (Jump to register) operation.
// It updates the R7 register with the current PC value and then jumps to the target address.
// The target address can be either a long address or a base register address, depending on the longFlag.
// If longFlag is 1, the target address is calculated by adding the sign-extended PC offset to the current PC.
// If longFlag is 0, the target address is the value stored in the base register.
func opJSR(instruction uint16) {
	var longFlag = (instruction >> 11) & 1
	Register[R_R7] = Register[R_PC]
	if longFlag == 1 {
		var pcOffset = signExtend(instruction&0x7ff, 11)
		Register[R_PC] += pcOffset
	} else {
		var baseR = (instruction >> 6) & 0x7
		Register[R_PC] = Register[baseR]
	}
}

// opAND performs the bitwise AND operation on the specified operands.
// It updates the value of the destination register based on the result of the operation.
// The operands can be either registers or immediate values, depending on the immFlag.
// If immFlag is set, the immediate value is sign-extended and bitwise ANDed with the value in register r1,
// and the result is stored in register r0.
// If immFlag is not set, the values in registers r1 and r2 are bitwise ANDed, and the result is stored in register r0.
// The function also updates the condition flags based on the result.
func opAND(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)
	var immFlag = getImmFlag(instruction)

	if immFlag == 1 {
		var imm5 = signExtend(instruction&0x1f, 5)
		Register[r0] = Register[r1] & imm5
	} else {
		var r2 = instruction & 0x7
		Register[r0] = Register[r1] & Register[r2]
	}
	updateFlags(r0)
}

// opLDR loads a value from memory into a register.
// It takes an instruction as input and extracts the necessary fields to perform the operation.
// The value is loaded from the memory address calculated by adding the offset to the value in the source register.
// The loaded value is then stored in the destination register.
// Finally, the flags are updated based on the result stored in the destination register.
func opLDR(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)
	var offset = signExtend(instruction&0x3f, 6)
	Register[r0] = memory.Read(Register[r1] + offset)
	updateFlags(r0)
}

// opSTR stores the value in register r0 into memory at the address calculated by adding the offset to the value in register r1.
// The offset is sign-extended to 16 bits before adding it to the base address.
func opSTR(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)
	var offset = signExtend(instruction&0x3f, 6)
	memory.Write(Register[r1]+offset, Register[r0])
}

func opRTI(instruction uint16) {
	log.Printf("Operation code not implemented: 0x%04X", instruction)
}

// opNOT performs a bitwise NOT operation on the value in the register specified by r1
// and stores the result in the register specified by r0.
// It also updates the flags based on the result.
func opNOT(instruction uint16) {
	var r0 = getR0(instruction)
	var r1 = getR1(instruction)
	Register[r0] = ^Register[r1]
	updateFlags(r0)
}

// opLDI loads a value from memory using an indirect address.
// It takes an instruction as input and updates the specified register with the loaded value.
// The instruction is expected to be a 16-bit unsigned integer.
func opLDI(instruction uint16) {
	var r0 = getR0(instruction)
	var pcOffset = signExtend(instruction&0x1ff, 9)
	Register[r0] = memory.Read(memory.Read(Register[R_PC] + pcOffset))
	updateFlags(r0)
}

// opSTI stores the value in the specified register into memory at the address
// calculated by adding the program counter (PC) and the sign-extended offset.
// The register to be stored is determined by the value of the first 3 bits of
// the instruction.
// Parameters:
//   - instruction: The 16-bit instruction to be executed.
func opSTI(instruction uint16) {
	var r0 = getR0(instruction)
	var pcOffset = signExtend(instruction&0x1ff, 9)
	memory.Write(memory.Read(Register[R_PC]+pcOffset), Register[r0])
}

// opJMP performs an unconditional jump to the address specified by the register r1.
func opJMP(instruction uint16) {
	var r1 = getR1(instruction)
	Register[R_PC] = Register[r1]
}

func opRES(instruction uint16) {
	log.Printf("Operation code not implemented: 0x%04X", instruction)
}

// opLEA loads the effective address of the specified PC offset into the register r0.
func opLEA(instruction uint16) {
	var r0 = getR0(instruction)
	var pcOffset = signExtend(instruction&0x1ff, 9)
	Register[r0] = Register[R_PC] + pcOffset
	updateFlags(r0)
}

func opTRAP(instruction uint16) {
	// Extract the trap vector from the instruction
}

// callOpcode calls the function associated with the given opcode and passes the instruction as an argument.
// If the opcode is not found in the OPCODES map, it logs an error message.
func callOpcode(opcode uint16, instruction uint16) {
	if f, ok := OPCODES[opcode]; ok {
		f(instruction)
	} else {
		log.Printf("invalid opcode=%v", opcode)
	}
}
