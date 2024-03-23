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

func opBR(instruction uint16) {
	// TODO
}

func opADD(instruction uint16) {
	// TODO
}

func opLD(instruction uint16) {
	// TODO
}

func opST(instruction uint16) {

}

func opJSR(instruction uint16) {

}

func opAND(instruction uint16) {

}

func opLDR(instruction uint16) {

}

func opSTR(instruction uint16) {

}

func opRTI(instruction uint16) {

}

func opNOT(instruction uint16) {

}

func opLDI(instruction uint16) {

}

func opSTI(instruction uint16) {

}

func opJMP(instruction uint16) {

}

func opRES(instruction uint16) {

}

func opLEA(instruction uint16) {

}

func opTRAP(instruction uint16) {

}

func callOpcode(opcode uint16, instruction uint16) {
	if f, ok := OPCODES[opcode]; ok {
		f(instruction)
	} else {
		log.Printf("invalid opcode=%v", opcode)
	}

}
