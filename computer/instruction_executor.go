package computer

/*
	Remember to check the pseudo-ops and psuedo-instructions!
*/

/*The InstructionExecutor is responsible for tracking the Current Instruction and
executing the instruction using an internal object.
*/
type InstructionExecutor struct {
	operator           instructionOperator
	instructionAddress uint16
	instruction        uint32
	memory             instructionReadWriteMemory
}

type instructionOperator interface {
	load(reg uint, address uint16, m instructionReadMemory)
	store(reg uint, address uint16, m instructionWriteMemory)
	add(dest uint, reg1 uint, reg2 uint)
	sub(dest uint, reg1 uint, reg2 uint)
	bit_and(dest uint, reg1 uint, reg2 uint)
	bit_or(dest uint, reg1 uint, reg2 uint)
	left_shift(dest uint, reg uint)
	right_shift(dest uint, reg uint, preserveSign bool)
	multiply(dest uint, reg1 uint, reg2 uint)
	divide(destDividend uint, destRem uint, reg1 uint, reg2 uint)
}

type instructionReadMemory interface {
	Get(address uint16) uint32
}

type instructionWriteMemory interface {
	Set(address uint16, value uint32)
}

type instructionReadWriteMemory interface {
	instructionReadMemory
	instructionWriteMemory
}

func (ex *InstructionExecutor) execute() {
	ex.loadInstruction()
	ex.executeInstruction()
}

func (ex *InstructionExecutor) loadInstruction() {
	ex.instruction = ex.memory.Get(ex.instructionAddress)
}

func (ex *InstructionExecutor) executeInstruction() {

}

func (ex *InstructionExecutor) addImmediate() {

}

func (ex *InstructionExecutor) setLessThanImmediate() {

}

func (ex *InstructionExecutor) setLessThanImmediateUnsigned() {

}

func (ex *InstructionExecutor) andImmmediate() {

}

func (ex *InstructionExecutor) orImmediate() {

}

func (ex *InstructionExecutor) xorImmediate() {

}

func (ex *InstructionExecutor) shiftLeftLogicalImmediate() {

}

func (ex *InstructionExecutor) shiftRightLogicalImmediate() {

}

func (ex *InstructionExecutor) shiftRightArithmeticImmediate() {

}

func (ex *InstructionExecutor) loadUpperImmediate() {

}

func (ex *InstructionExecutor) addUpperImmediateToPC() {

}

func (ex *InstructionExecutor) add() {

}

func (ex *InstructionExecutor) setLessThan() {

}

func (ex *InstructionExecutor) setLessThanUnsigned() {

}

func (ex *InstructionExecutor) and() {

}

func (ex *InstructionExecutor) or() {

}

func (ex *InstructionExecutor) xor() {

}

func (ex *InstructionExecutor) shiftLeftLogical() {

}

func (ex *InstructionExecutor) shiftRightLogical() {

}

func (ex *InstructionExecutor) shiftRightArithmetic() {

}

func (ex *InstructionExecutor) nop() {

}

func (ex *InstructionExecutor) branchEqual() {

}

func (ex *InstructionExecutor) branchNotEqual() {

}

func (ex *InstructionExecutor) branchLessThan() {

}

func (ex *InstructionExecutor) branchLessThanUnsigned() {

}

func (ex *InstructionExecutor) branchGreaterThanOrEqual() {

}

func (ex *InstructionExecutor) branchGreaterThanOrEqualUnsigned() {

}

func (ex *InstructionExecutor) jumpAndLink() {

}

func (ex *InstructionExecutor) jumpAndLinkRegister() {

}

func (ex *InstructionExecutor) loadWord() {

}

func (ex *InstructionExecutor) loadHalfWord() {

}

func (ex *InstructionExecutor) loadHalfWordUnsigned() {

}

func (ex *InstructionExecutor) loadByte() {

}

func (ex *InstructionExecutor) loadByteUnsigned() {

}

func (ex *InstructionExecutor) storeWord() {

}

func (ex *InstructionExecutor) storeHalfWord() {

}

func (ex *InstructionExecutor) storeByte() {

}

/* CONTINUE FROM 2.7 FOR ADDITIONAL INSTRUCTIONS */
