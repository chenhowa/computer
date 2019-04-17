package computer

import (
	"math"
)

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
	add_immediate(dest uint, reg uint, immediate uint32)
	sub(dest uint, reg1 uint, reg2 uint)
	bit_and(dest uint, reg1 uint, reg2 uint)
	bit_and_immediate(dest uint, reg uint, immediate uint32)
	bit_or(dest uint, reg1 uint, reg2 uint)
	bit_or_immediate(dest uint, reg uint, immediate uint32)
	bit_xor(dest uint, reg1 uint, reg2 uint)
	bit_xor_immediate(dest uint, reg uint, immediate uint32)
	left_shift(dest uint, reg uint)
	right_shift(dest uint, reg uint, preserveSign bool)
	multiply(dest uint, reg1 uint, reg2 uint)
	divide(destDividend uint, destRem uint, reg1 uint, reg2 uint)
	get(reg uint) uint32
	zero(dest uint)
	one(dest uint)
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

func (ex *InstructionExecutor) addImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Add and ignore overflow.
	ex.operator.add_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *InstructionExecutor) setLessThanImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Compare as signed numbers
	// 3. Place 1 or 0 in destination reg, based on result.

	regValueLess := int32(ex.operator.get(reg)) < int32(signExtendUint32WithBit(immediate, 11))
	if regValueLess {
		ex.operator.one(dest)
	} else {
		ex.operator.zero(dest)
	}
}

func (ex *InstructionExecutor) setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32) {
	regValueLess := (ex.operator.get(reg)) < (signExtendUint32WithBit(immediate, 11))
	if regValueLess {
		ex.operator.one(dest)
	} else {
		ex.operator.zero(dest)
	}
}

func (ex *InstructionExecutor) andImmmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_and_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *InstructionExecutor) orImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_or_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *InstructionExecutor) xorImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_xor_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
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

func (ex *InstructionExecutor) csrReadAndWrite() {

}

func (ex *InstructionExecutor) csrReadAndSet() {

}

func (ex *InstructionExecutor) csrReadAndClear() {

}

func (ex *InstructionExecutor) csrReadAndWriteImmediate() {

}

func (ex *InstructionExecutor) csrReadAndSetImmediate() {

}

func (ex *InstructionExecutor) csrReadAndClearImmediate() {

}

func (ex *InstructionExecutor) envCall() {

}

func (ex *InstructionExecutor) envBreak() {

}
