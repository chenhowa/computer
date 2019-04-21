package computer

import (
	"math"
)

/*
	Remember to check the pseudo-ops and psuedo-instructions!
*/

/*The RiscVInstructionExecutor is responsible for translating
a binary instruction into MESSAGES to the Operator and InstructionManager
executing the instruction using an internal object.
*/
type RiscVInstructionExecutor struct {
	operator instructionOperator
	// A translator doesn't ned to know about memory...
	//memory             instructionReadWriteMemory
	instructionManager instructionManager
}

type instructionManager interface {
	getInstructionAddress() uint32
	getInstruction(memory instructionReadMemory) uint32
	incrementInstructionAddress()
	addOffset(offset uint32)
	loadInstructionAddress(newAddress uint32)
}

type instructionOperator interface {
	load_word(dest uint, address uint32, memory instructionReadMemory)
	load_halfword(dest uint, address uint32, memory instructionReadMemory)
	load_halfword_unsigned(dest uint, address uint32, memory instructionReadMemory)
	load_byte(dest uint, address uint32, memory instructionReadMemory)
	load_byte_unsigned(dest uint, address uint32, memory instructionReadMemory)
	store_word(src uint, address uint32, memory instructionWriteMemory)
	store_halfword(src uint, address uint32, memory instructionWriteMemory)
	store_byte(src uint, address uint32, memory instructionWriteMemory)
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
	Get(address uint32) uint32
}

type instructionWriteMemory interface {
	Set(address uint32, value uint32)
}

type instructionReadWriteMemory interface {
	instructionReadMemory
	instructionWriteMemory
}

/*
func (ex *RiscVInstructionExecutor) execute(memory instructionReadWriteMemory) {
	instruction := ex.instructionManager.getInstruction(memory)
	ex.instructionManager.incrementInstructionAddress()
	ex.executeInstruction(instruction, memory)
}

func (ex *RiscVInstructionExecutor) executeInstruction(instruction uint32, memory instructionReadWriteMemory) {
	blah
}*/

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

func (ex *RiscVInstructionExecutor) AddImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Add and ignore overflow.
	ex.operator.add_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) SetLessThanImmediate(dest uint, reg uint, immediate uint32) {
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Compare as signed numbers
	// 3. Place 1 or 0 in destination reg, based on result.

	regValueLess := int32(ex.operator.get(reg)) < int32(signExtendUint32WithBit(immediate, 11))
	ex.operator.bit_and_immediate(dest, dest, 0) // zero out destination
	if regValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) SetLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32) {
	regValueLess := (ex.operator.get(reg)) < (signExtendUint32WithBit(immediate, 11))
	ex.operator.bit_and_immediate(dest, dest, 0)
	if regValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) AndImmmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_and_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) OrImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_or_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) XorImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.bit_xor_immediate(dest, reg, signExtendUint32WithBit(immediate, 11))
}

func (ex *RiscVInstructionExecutor) ShiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.left_shift_immediate(dest, reg, immediate)
}

func (ex *RiscVInstructionExecutor) ShiftRightLogicalImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.right_shift_immediate(dest, reg, immediate, false)
}

func (ex *RiscVInstructionExecutor) ShiftRightArithmeticImmediate(dest uint, reg uint, immediate uint32) {
	ex.operator.right_shift_immediate(dest, reg, immediate, true)
}

/*
	Takes the lower 20 bits of the immediate and loads them into the
	upper 20 bits of the destination register, with the lower 20 bits
	of the destination register as 0
*/
func (ex *RiscVInstructionExecutor) LoadUpperImmediate(dest uint, immediate uint32) {
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, immediate<<12)
}

func (ex *RiscVInstructionExecutor) AddUpperImmediateToPC(dest uint, immediate uint32) {
	result := (immediate << 12) + uint32(ex.instructionManager.getInstructionAddress())
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, result)
}

func (ex *RiscVInstructionExecutor) Add(dest uint, reg1 uint, reg2 uint) {
	ex.operator.add(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) SetLessThan(dest uint, reg1 uint, reg2 uint) {
	reg1ValueLess := int32(ex.operator.get(reg1)) < int32(ex.operator.get(reg2))
	ex.operator.bit_and_immediate(dest, dest, 0) // zero out destination
	if reg1ValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) SetLessThanUnsigned(dest uint, reg1 uint, reg2 uint) {
	reg1ValueLess := (ex.operator.get(reg1)) < ex.operator.get(reg2)
	ex.operator.bit_and_immediate(dest, dest, 0)
	if reg1ValueLess {
		ex.operator.bit_or_immediate(dest, dest, 1)
	}
}

func (ex *RiscVInstructionExecutor) And(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_and(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) Or(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_or(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) Xor(dest uint, reg1 uint, reg2 uint) {
	ex.operator.bit_xor(dest, reg1, reg2)
}

func (ex *RiscVInstructionExecutor) ShiftLeftLogical(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (math.MaxUint32 >> (32 - 5))
	ex.operator.left_shift_immediate(dest, reg, lowerFiveBits)
}

func (ex *RiscVInstructionExecutor) ShiftRightLogical(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (uint32(math.MaxUint32) >> (32 - 5))
	ex.operator.right_shift_immediate(dest, reg, lowerFiveBits, false)
}

func (ex *RiscVInstructionExecutor) ShiftRightArithmetic(dest uint, reg uint, shiftreg uint) {
	lowerFiveBits := ex.operator.get(shiftreg) & (math.MaxUint32 >> (32 - 5))
	ex.operator.right_shift_immediate(dest, reg, lowerFiveBits, true)
}

func (ex *RiscVInstructionExecutor) Nop() {
	// This instruction literally does nothing.
	// In RISC-V, it can be encoded as ADDI x0, x0, 0,
	// but there is no need for that here.
}

func (ex *RiscVInstructionExecutor) BranchEqual() {

}

func (ex *RiscVInstructionExecutor) BranchNotEqual() {

}

func (ex *RiscVInstructionExecutor) BranchLessThan() {

}

func (ex *RiscVInstructionExecutor) BranchLessThanUnsigned() {

}

func (ex *RiscVInstructionExecutor) BranchGreaterThanOrEqual() {

}

func (ex *RiscVInstructionExecutor) BranchGreaterThanOrEqualUnsigned() {

}

/*
	This instruction saves the address of the next instruction into the
	destination register and then adds the lower 20 bits of
	the offset to the program counter (stored in the program counter)
*/
func (ex *RiscVInstructionExecutor) JumpAndLink(dest uint, pcOffset uint32) {
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, ex.instructionManager.getInstructionAddress())

	lower20Bits := pcOffset & (uint32(math.MaxUint32) >> (32 - 20))
	ex.instructionManager.addOffset(lower20Bits)
}

/*
	This instruction saves the address of the next instruction into the
	destination register. It then adds the lower 12 bits of the offset
	to the value in the Base Register, sets the LSB of the result to 0,
	It then sets the program counter to the result.
*/
func (ex *RiscVInstructionExecutor) JumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32) {
	ex.operator.bit_and_immediate(dest, dest, 0)
	ex.operator.bit_or_immediate(dest, dest, ex.instructionManager.getInstructionAddress())

	lower12Bits := pcOffset & (uint32(math.MaxUint32) >> (32 - 12))
	address := lower12Bits + ex.operator.get(basereg)
	ex.instructionManager.loadInstructionAddress(address)
}

func (ex *RiscVInstructionExecutor) LoadWord(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	lower12Bits := offset & (uint32(math.MaxUint32) >> (32 - 12))
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.load_word(dest, address, memory)
}

func (ex *RiscVInstructionExecutor) LoadHalfWord(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.load_halfword(dest, address, memory)

}

func (ex *RiscVInstructionExecutor) LoadHalfWordUnsigned(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.load_halfword_unsigned(dest, address, memory)

}

func (ex *RiscVInstructionExecutor) LoadByte(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.load_byte(dest, address, memory)

}

func (ex *RiscVInstructionExecutor) LoadByteUnsigned(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.load_byte_unsigned(dest, address, memory)

}

func (ex *RiscVInstructionExecutor) StoreWord(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.store_word(src, address, memory)
}

func (ex *RiscVInstructionExecutor) StoreHalfWord(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.store_halfword(src, address, memory)
}

func (ex *RiscVInstructionExecutor) StoreByte(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	lower12Bits := offset & uint32(math.MaxUint32) >> (32 - 12)
	address := lower12Bits + ex.operator.get(reg)
	ex.operator.store_byte(src, address, memory)
}

func (ex *RiscVInstructionExecutor) CsrReadAndWrite() {

}

func (ex *RiscVInstructionExecutor) CsrReadAndSet() {

}

func (ex *RiscVInstructionExecutor) CsrReadAndClear() {

}

func (ex *RiscVInstructionExecutor) CsrReadAndWriteImmediate() {

}

func (ex *RiscVInstructionExecutor) CsrReadAndSetImmediate() {

}

func (ex *RiscVInstructionExecutor) CsrReadAndClearImmediate() {

}

func (ex *RiscVInstructionExecutor) EnvCall() {

}

func (ex *RiscVInstructionExecutor) EnvBreak() {

}
