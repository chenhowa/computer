package computer

import (
	"math"
)

/*
	Remember to check the pseudo-ops and psuedo-instructions!
*/

/*The RiscVInstructionExecutor is responsible for tracking the Current Instruction and
executing the instruction using an internal object.
*/
type RiscVInstructionExecutor struct {
	operator           instructionOperator
	instructionAddress uint32
	instruction        uint32
	memory             instructionReadWriteMemory
}

type instructionOperator interface {
	load(reg uint, address uint16, m instructionReadMemory)
	store(reg uint, address uint16, m instructionWriteMemory)
	add(dest uint, reg1 uint, reg2 uint)
	add_immediate(dest uint, reg uint, immediate uint32)
	sub(dest uint, reg1 uint, reg2 uint)
	bit_and(dest uint, reg1 uint, reg2 uint)
	bit_and_immediate(dest uint, reg uint, immediate uint32)
	bit_or(dest uint, reg1 uint, reg2 uint)
	bit_or_immediate(dest uint, reg uint, immediate uint32)
	bit_xor(dest uint, reg1 uint, reg2 uint)
	bit_xor_immediate(dest uint, reg uint, immediate uint32)
	left_shift_immediate(dest uint, reg uint, immediate uint32)
	right_shift_immediate(dest uint, reg uint, immediate uint32, preserveSign bool)
	multiply(dest uint, reg1 uint, reg2 uint)
	divide(destDividend uint, destRem uint, reg1 uint, reg2 uint)
	get(reg uint) uint32
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

func (ex *RiscVInstructionExecutor) execute() {
	ex.loadInstruction()
	ex.executeInstruction()
}

func (ex *RiscVInstructionExecutor) loadInstruction() {
	// does the type conversion truncate the top or bottom
	// bits?
	ex.instruction = ex.memory.Get(uint16(ex.instructionAddress))
}

func (ex *RiscVInstructionExecutor) executeInstruction() {

}

func signExtendUint32WithBit(integer uint32, bit uint) uint32 {
	bitValue := integer >> (bit - 1)
	var mask uint32
	var signExtended uint32
	if max := math.MaxUint32; bitValue == 1 {
		mask = uint32(max << (bit + 1))
		signExtended = mask | integer
	} else {
		mask = uint32(max >> (32 - bit - 1))
		signExtended = mask & integer
	}

	return signExtended
}

func (ex *RiscVInstructionExecutor) addImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Add and ignore overflow.
	ex.operator.add_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) setLessThanImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Compare as signed numbers
	// 3. Place 1 or 0 in destination reg, based on result.

	regValueLess := int32(ex.operator.get(reg)) < int32(signExtendUint32WithBit(immediate, 11))
	ex.operator.bit_and_immediate(dest, dest, 0) // zero out destination
	if regValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32) {
	regValueLess := (ex.operator.get(reg)) < (signExtendUint32WithBit(immediate, 11))
	ex.operator.bit_and_immediate(dest, dest, 0)
	if regValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) andImmmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_and_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) orImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_or_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) xorImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_xor_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) shiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.left_shift_immediate(dest, reg, immediate)
}

func (ex *RiscVInstructionExecutor) shiftRightLogicalImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.right_shift_immediate(dest, reg, immediate, false)
}

func (ex *RiscVInstructionExecutor) shiftRightArithmeticImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.right_shift_immediate(dest, reg, immediate, true)
}

/*
	Takes the lower 20 bits of the immediate and loads them into the
	upper 20 bits of the destination register, with the lower 20 bits
	of the destination register as 0
*/
func (ex *RiscVInstructionExecutor) loadUpperImmediate(dest uint, immediate uint32) {
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, immediate<<12)
}

func (ex *RiscVInstructionExecutor) addUpperImmediateToPC(dest uint, immediate uint32) {
	result := (immediate << 12) + uint32(ex.instructionAddress)
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, result)
}

func (ex *RiscVInstructionExecutor) add(dest uint, reg1 uint, reg2 uint) {
	ex.operator.add(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) setLessThan(dest uint, reg1 uint, reg2 uint) {
	reg1ValueLess := int32(ex.operator.get(reg1)) < int32(ex.operator.get(reg2))
	ex.operator.bit_and_immediate(dest, dest, 0) // zero out destination
	if reg1ValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) setLessThanUnsigned(dest uint, reg1 uint, reg2 uint) {
	reg1ValueLess := (ex.operator.get(reg1)) < ex.operator.get(reg2)
	ex.operator.bit_and_immediate(dest, dest, 0)
	if reg1ValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) and(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_and(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) or(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_or(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) xor(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_xor(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) shiftLeftLogical(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (math.MaxUint32 >> (32 - 5))
	ex.operator.left_shift_immediate(dest, reg, lowerFiveBits)
}

func (ex *RiscVInstructionExecutor) shiftRightLogical(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (uint32(math.MaxUint32) >> (32 - 5))
	ex.operator.right_shift_immediate(dest, reg, lowerFiveBits, false)
}

func (ex *RiscVInstructionExecutor) shiftRightArithmetic(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (math.MaxUint32 >> (32 - 5))
	ex.operator.right_shift_immediate(dest, reg, lowerFiveBits, true)
}

func (ex *RiscVInstructionExecutor) nop() {
	// This instruction literally does nothing.
	// In RISC-V, it can be encoded as ADDI x0, x0, 0,
	// but there is no need for that here.
}

func (ex *RiscVInstructionExecutor) branchEqual() {

}

func (ex *RiscVInstructionExecutor) branchNotEqual() {

}

func (ex *RiscVInstructionExecutor) branchLessThan() {

}

func (ex *RiscVInstructionExecutor) branchLessThanUnsigned() {

}

func (ex *RiscVInstructionExecutor) branchGreaterThanOrEqual() {

}

func (ex *RiscVInstructionExecutor) branchGreaterThanOrEqualUnsigned() {

}

/*
	This instruction saves (program counter + 4) into the
	destination register and then adds the lower 20 bits of
	the offset to the program counter (stored in the program counter)
*/
func (ex *RiscVInstructionExecutor) jumpAndLink(dest uint, pcOffset uint32) {
	ex.operator.bit_or_immediate(dest, dest, 0)
	ex.operator.bit_and_immediate(dest, dest, ex.instructionAddress+4)

	lower20Bits := pcOffset & uint32(math.MaxUint32) >> (32 - 20)
	ex.instructionAddress += lower20Bits
}

/*
	This instruction saves (program counter + 4) into the
	destination register. It then adds the lower 12 bits of the offset
	to the value in the Base Register, sets the LSB of the result to 0,
	It then sets the program counter to the result.
*/
func (ex *RiscVInstructionExecutor) jumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32) {
	ex.operator.bit_or_immediate(dest, dest, 0)
	ex.operator.bit_and_immediate(dest, dest, ex.instructionAddress+4)

	lower12Bits := pcOffset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(basereg)
	ex.instructionAddress = address
}

func (ex *RiscVInstructionExecutor) loadWord() {

}

func (ex *RiscVInstructionExecutor) loadHalfWord() {

}

func (ex *RiscVInstructionExecutor) loadHalfWordUnsigned() {

}

func (ex *RiscVInstructionExecutor) loadByte() {

}

func (ex *RiscVInstructionExecutor) loadByteUnsigned() {

}

func (ex *RiscVInstructionExecutor) storeWord() {

}

func (ex *RiscVInstructionExecutor) storeHalfWord() {

}

func (ex *RiscVInstructionExecutor) storeByte() {

}

func (ex *RiscVInstructionExecutor) csrReadAndWrite() {

}

func (ex *RiscVInstructionExecutor) csrReadAndSet() {

}

func (ex *RiscVInstructionExecutor) csrReadAndClear() {

}

func (ex *RiscVInstructionExecutor) csrReadAndWriteImmediate() {

}

func (ex *RiscVInstructionExecutor) csrReadAndSetImmediate() {

}

func (ex *RiscVInstructionExecutor) csrReadAndClearImmediate() {

}

func (ex *RiscVInstructionExecutor) envCall() {

}

func (ex *RiscVInstructionExecutor) envBreak() {

}
