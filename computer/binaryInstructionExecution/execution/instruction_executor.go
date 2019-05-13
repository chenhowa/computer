package execution

import (
	"math"

	Utils "github.com/chenhowa/os/computer/binaryInstructionExecution/bitUtils"
)

/*The RiscVInstructionExecutor is responsible for taking the operands of
a binary instruction and turning them into MESSAGES to the Operator and InstructionManager
executing the instruction using an internal object.
*/
type RiscVInstructionExecutor struct {
	operator instructionOperator
}

type executionEnvManager interface {
	executeCall()
}

type debugEnvManager interface {
	debugBreak()
}

type csrOperator interface {
	get(reg uint) uint32
	set(reg uint, val uint32)
}

type instructionManager interface {
	getCurrentInstructionAddress() uint32
	getNextInstructionAddress() uint32
	addOffsetForNextInstructionAddress(offset uint32)
	loadAsNextInstructionAddress(newAddress uint32)
}

type instructionOperator interface {
	loadWord(dest uint, address uint32, memory instructionReadMemory)
	loadHalfWord(dest uint, address uint32, memory instructionReadMemory)
	loadHalfWordUnsigned(dest uint, address uint32, memory instructionReadMemory)
	loadByte(dest uint, address uint32, memory instructionReadMemory)
	loadByteUnsigned(dest uint, address uint32, memory instructionReadMemory)
	storeWord(src uint, address uint32, memory instructionWriteMemory)
	storeHalfWord(src uint, address uint32, memory instructionWriteMemory)
	storeByte(src uint, address uint32, memory instructionWriteMemory)
	add(dest uint, reg1 uint, reg2 uint)
	addImmediate(dest uint, reg uint, immediate uint32)
	sub(dest uint, reg1 uint, reg2 uint)
	and(dest uint, reg1 uint, reg2 uint)
	andImmediate(dest uint, reg uint, immediate uint32)
	or(dest uint, reg1 uint, reg2 uint)
	orImmediate(dest uint, reg uint, immediate uint32)
	xor(dest uint, reg1 uint, reg2 uint)
	xorImmediate(dest uint, reg uint, immediate uint32)
	leftShiftImmediate(dest uint, reg uint, immediate uint32)
	rightShiftImmediate(dest uint, reg uint, immediate uint32, preserveSign bool)
	multiply(dest uint, reg1 uint, reg2 uint)
	divide(destDividend uint, destRem uint, reg1 uint, reg2 uint)
	get(reg uint) uint32
}

type instructionReadMemory interface {
	Get(address uint32) uint32
}

type instructionWriteMemory interface {
	Set(address uint32, value uint32, bitsToSet uint)
}

type instructionReadWriteMemory interface {
	instructionReadMemory
	instructionWriteMemory
}

/* resetRegisterZero resets the value of register 0 to 0,
since according to RiscV, the value of register 0 should always be 0 */
func (ex *RiscVInstructionExecutor) resetRegisterZero() {
	ex.operator.andImmediate(0, 0, 0)
}

/*Get allows caller to read the value of register `reg`. In RiscV, any of the 32 registers can be read,
but the 0 register cannot be written to -- the value of the 0 register is always 0
*/
func (ex *RiscVInstructionExecutor) Get(reg uint) uint32 {
	defer ex.resetRegisterZero()

	if reg == 0 {
		return 0
	}

	return ex.operator.get(reg)
}

/*AddImmediate adds an immediate to a value in a register, and stores the result
in the destination register. The immediate value is `immediate` 12 least-significant
bits, sign-extended based on the 12th bit.
*/
func (ex *RiscVInstructionExecutor) AddImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()

	// 1. Sign extend the immediate based on the 11th bit
	// 2. Add and ignore overflow.
	ex.operator.addImmediate(dest, reg, Utils.SignExtendUint32WithBit(immediate, 11))
}

/*SetLessThanImmediate compares the value in a register to the sign-extended immediate (using the 11th bit
to sign extend), and sets 1 if the register value is less, otherwise it sets 0
*/
func (ex *RiscVInstructionExecutor) SetLessThanImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	// 1. Sign extend the immediate based on the 12th bit
	// 2. Compare as signed numbers
	// 3. Place 1 or 0 in destination reg, based on result.

	regValueLess := int32(ex.Get(reg)) < int32(Utils.SignExtendUint32WithBit(immediate, 11))
	ex.operator.andImmediate(dest, dest, 0) // zero out destination
	if regValueLess {
		ex.operator.orImmediate(dest, dest, 1)
	}
}

/*SetLessThanImmediateUnsigned compares the value in a register to the unsigned immediate,
and sets 1 if the register value is less, otherwise it sets 0
*/
func (ex *RiscVInstructionExecutor) SetLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	regValueLess := (ex.Get(reg)) < (Utils.SignExtendUint32WithBit(immediate, 11))
	ex.operator.andImmediate(dest, dest, 0)
	if regValueLess {
		ex.operator.orImmediate(dest, dest, 1)
	}
}

/*AndImmediate takes the `immediate`, sign-extends it with 12th bit, does a bit-wise AND with
the value in the register `reg`, and places the result in the register `dest` */
func (ex *RiscVInstructionExecutor) AndImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	ex.operator.andImmediate(dest, reg, Utils.SignExtendUint32WithBit(immediate, 11))
}

/*OrImmediate takes the `immediate`, sign-extends it with 12th bit, does a bit-wise OR with
the value in the register `reg`, and places the result in the register `dest` */
func (ex *RiscVInstructionExecutor) OrImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	ex.operator.orImmediate(dest, reg, Utils.SignExtendUint32WithBit(immediate, 11))
}

/*XorImmediate takes the `immediate`, sign-extends it with 12th bit, does a bit-wise XOR with
the value in the register `reg`, and places the result in the register `dest` */
func (ex *RiscVInstructionExecutor) XorImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	ex.operator.xorImmediate(dest, reg, Utils.SignExtendUint32WithBit(immediate, 11))
}

/*ShiftLeftLogicalImmediate takes the lowest 5 bits of `immediate`, and left-shifts the value
in register `reg` by that 5-bit amount, and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	shiftAmt := Utils.KeepBitsInInclusiveRange(immediate, 0, 4)
	ex.operator.leftShiftImmediate(dest, reg, shiftAmt)
}

/*ShiftRightLogicalImmediate takes the lowest 5 bits of `immediate`, and logical right-shifts the value
in register `reg` by that 5-bit amount, and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftRightLogicalImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	shiftAmt := Utils.KeepBitsInInclusiveRange(immediate, 0, 4)
	ex.operator.rightShiftImmediate(dest, reg, shiftAmt, false)
}

/*ShiftRightArithmeticImmediate takes the lowest 5 bits of `immediate`, and arithmetic right-shifts the value
in register `reg` by that 5-bit amount, and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftRightArithmeticImmediate(dest uint, reg uint, immediate uint32) {
	defer ex.resetRegisterZero()
	shiftAmt := Utils.KeepBitsInInclusiveRange(immediate, 0, 4)
	ex.operator.rightShiftImmediate(dest, reg, shiftAmt, true)
}

/*LoadUpperImmediate takes the lower 20 bits of `immediate` and loads them into the
upper 20 bits of the register `dest`, with the lower 12 bits of `dest` set to 0
*/
func (ex *RiscVInstructionExecutor) LoadUpperImmediate(dest uint, immediate uint32) {
	defer ex.resetRegisterZero()
	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, immediate<<12)
}

/*AddUpperImmediateToPC takes the lower 20 bits of `immediate`, left-shifts by 12 bits, and then adds this
offset to the value in the Program Counter (PC), and places the result in the register `dest`*/
func (ex *RiscVInstructionExecutor) AddUpperImmediateToPC(dest uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	result := (immediate << 12) + uint32(manager.getCurrentInstructionAddress())
	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, result)
}

/*Add take the values in registers `reg1` and `reg2`, adds them together (ignoring overflow), and writes
the lower 32 bits into register `dest`
*/
func (ex *RiscVInstructionExecutor) Add(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	ex.operator.add(dest, reg1, reg2)
}

/*Sub subtracts the value in register `reg2` from the value in `reg1`, ignoring overflow,
and writes the result's lower 32 bits into register `dest` */
func (ex *RiscVInstructionExecutor) Sub(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	ex.operator.sub(dest, reg1, reg2)
}

/*SetLessThan performs a signed comparison of the values in registers `reg1` and `reg2`.
If `reg1` is les than `reg2`, register `dest` is set to 1; otherwise it is set to 0*/
func (ex *RiscVInstructionExecutor) SetLessThan(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	reg1ValueLess := int32(ex.Get(reg1)) < int32(ex.Get(reg2))
	ex.operator.andImmediate(dest, dest, 0) // zero out destination
	if reg1ValueLess {
		ex.operator.orImmediate(dest, dest, 1)
	}
}

/*SetLessThanUnsigned performs an unsigned comparison of the values in registers `reg1` and `reg2`.
If `reg1` is les than `reg2`, register `dest` is set to 1; otherwise it is set to 0*/
func (ex *RiscVInstructionExecutor) SetLessThanUnsigned(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	reg1ValueLess := (ex.Get(reg1)) < ex.Get(reg2)
	ex.operator.andImmediate(dest, dest, 0)
	if reg1ValueLess {
		ex.operator.orImmediate(dest, dest, 1)
	}
}

/*And performs a bitwise AND between the values in registers `reg1` and `reg2`,
and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) And(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	ex.operator.and(dest, reg1, reg2)
}

/*Or performs a bitwise OR between the values in registers `reg1` and `reg2`,
and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) Or(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	ex.operator.or(dest, reg1, reg2)
}

/*Xor performs a bitwise Xor between the values in registers `reg1` and `reg2`,
and places the result in register `dest`*/
func (ex *RiscVInstructionExecutor) Xor(dest uint, reg1 uint, reg2 uint) {
	defer ex.resetRegisterZero()
	ex.operator.xor(dest, reg1, reg2)
}

/*ShiftLeftLogical takes the lowest 5 bits of the value in register `shiftreg`, and then
logical left-shifts the value in register `reg` by that amount, and places the result in
register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftLeftLogical(dest uint, reg uint, shiftreg uint) {
	defer ex.resetRegisterZero()
	lowerFiveBits := Utils.KeepBitsInInclusiveRange(ex.Get(shiftreg), 0, 4)
	ex.operator.leftShiftImmediate(dest, reg, lowerFiveBits)
}

/*ShiftRightLogical takes the lowest 5 bits of the value in register `shiftreg`, and then
logical right-shifts the value in register `reg` by that amount, and places the result in
register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftRightLogical(dest uint, reg uint, shiftreg uint) {
	defer ex.resetRegisterZero()
	lowerFiveBits := ex.Get(shiftreg) & (uint32(math.MaxUint32) >> (32 - 5))
	ex.operator.rightShiftImmediate(dest, reg, lowerFiveBits, false)
}

/*ShiftRightArithmetic takes the lowest 5 bits of the value in register `shiftreg`, and then
arithmetic right-shifts the value in register `reg` by that amount, and places the result in
register `dest`*/
func (ex *RiscVInstructionExecutor) ShiftRightArithmetic(dest uint, reg uint, shiftreg uint) {
	defer ex.resetRegisterZero()
	lowerFiveBits := ex.Get(shiftreg) & (math.MaxUint32 >> (32 - 5))
	ex.operator.rightShiftImmediate(dest, reg, lowerFiveBits, true)
}

/*BranchEqual compares the values in registers `reg1` and `reg2`. If `reg1` equals `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchEqual(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()

	if ex.Get(reg1) == ex.Get(reg2) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*BranchNotEqual compares the values in registers `reg1` and `reg2`. If `reg1` does NOT equal `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchNotEqual(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	if ex.Get(reg1) != ex.Get(reg2) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*BranchLessThan compares the values in registers `reg1` and `reg2` as SIGNED values. If `reg1` < `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchLessThan(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	if int32(ex.Get(reg1)) < int32(ex.Get(reg2)) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*BranchLessThanUnsigned compares the values in registers `reg1` and `reg2` as UNSIGNED values. If `reg1` < `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchLessThanUnsigned(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	if ex.Get(reg1) < ex.Get(reg2) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*BranchGreaterThanOrEqual compares the values in registers `reg1` and `reg2` as SIGNED values. If `reg1` >= `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchGreaterThanOrEqual(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	if int32(ex.Get(reg1)) >= int32(ex.Get(reg2)) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*BranchGreaterThanOrEqualUnsigned compares the values in registers `reg1` and `reg2` as UNSIGNED values. If `reg1` >= `reg2`, then
the 12 lowest bits of `immediate` are added to the pc through the `manager`
*/
func (ex *RiscVInstructionExecutor) BranchGreaterThanOrEqualUnsigned(reg1 uint, reg2 uint, immediate uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	if ex.Get(reg1) >= ex.Get(reg2) {
		manager.addOffsetForNextInstructionAddress(Utils.KeepBitsInInclusiveRange(immediate, 0, 11))
	}
}

/*JumpAndLink saves the address of the next instruction into the
register `dest` and then adds the sign-extended lowest 20 bits of
the `pcOffset` to the program counter using `manager`
*/
func (ex *RiscVInstructionExecutor) JumpAndLink(dest uint, pcOffset uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, manager.getNextInstructionAddress())

	lower20Bits := Utils.SignExtendUint32WithBit(Utils.KeepBitsInInclusiveRange(pcOffset, 0, 19), 19)
	manager.addOffsetForNextInstructionAddress(lower20Bits)
}

/*JumpAndLinkRegister saves the address of the next instruction into the
register `dest`. It then adds the sign-extended lowest 12 bits of `pcOffset`
to the value in the register `basereg`, sets the LSB of the result to 0,
and loads the result into the program counter through the `manager`
*/
func (ex *RiscVInstructionExecutor) JumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32, manager instructionManager) {
	defer ex.resetRegisterZero()
	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, manager.getNextInstructionAddress())

	lower12Bits := Utils.SignExtendUint32WithBit(Utils.KeepBitsInInclusiveRange(pcOffset, 0, 11), 11)
	address := Utils.KeepBitsInInclusiveRange(lower12Bits+ex.Get(basereg), 1, 31)
	manager.loadAsNextInstructionAddress(address)
}

/*LoadWord compiles an address from sign-extended lower 12 bits of offset, adds that to uint32 stored in
register `reg`, and then reads 1 Word from that address from memory into the destination register `dest`
*/
func (ex *RiscVInstructionExecutor) LoadWord(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.loadWord(dest, address, memory)
}

/*LoadHalfWord compiles an address from sign-extended lower 12 bits of offset, adds that to uint32 stored in
register `reg`, and then reads 1 Sign-extended HalfWord from that address from memory into the destination register `dest`
*/
func (ex *RiscVInstructionExecutor) LoadHalfWord(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.loadHalfWord(dest, address, memory)

}

/*LoadHalfWordUnsigned compiles an address from sign-extended lower 12 bits of offset, adds that to uint32 stored in
register `reg`, and then reads 1 HalfWord from that address from memory into the destination register `dest`
*/
func (ex *RiscVInstructionExecutor) LoadHalfWordUnsigned(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.loadHalfWordUnsigned(dest, address, memory)

}

/*LoadByte compiles an address from sign-extended lower 12 bits of offset, adds that to uint32 stored in
register `reg`, and then reads 1 Sign-extended Byte from that address from memory into the destination register `dest`
*/
func (ex *RiscVInstructionExecutor) LoadByte(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.loadByte(dest, address, memory)

}

/*LoadByteUnsigned compiles an address from sign-extended lower 12 bits of offset, adds that to uint32 stored in
register `reg`, and then reads 1 Byte from that address from memory into the destination register `dest`
*/
func (ex *RiscVInstructionExecutor) LoadByteUnsigned(dest uint, reg uint, offset uint32, memory instructionReadMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.loadByteUnsigned(dest, address, memory)

}

/*StoreWord compiles an address from the sign-extended lower 12 bits of `offset`, adds that to the uint32 stored
in register `reg`, and then stores 4 Bytes from the register `src` into memory at the calculated 32-bit address*/
func (ex *RiscVInstructionExecutor) StoreWord(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.storeWord(src, address, memory)
}

/*StoreHalfWord compiles an address from the sign-extended lower 12 bits of `offset`, adds that to the uint32 stored
in register `reg`, and then stores 2 Bytes from the register `src` into memory at the calculated 32-bit address*/
func (ex *RiscVInstructionExecutor) StoreHalfWord(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.storeHalfWord(src, address, memory)
}

/*StoreByte compiles an address from the sign-extended lower 12 bits of `offset`, adds that to the uint32 stored
in register `reg`, and then stores 1 Byte from the register `src` into memory at the calculated 32-bit address*/
func (ex *RiscVInstructionExecutor) StoreByte(src uint, reg uint, offset uint32, memory instructionWriteMemory) {
	defer ex.resetRegisterZero()
	lower12Bits := offset & Utils.KeepBitsInInclusiveRange(uint32(math.MaxUint32), 0, 11)
	address := Utils.SignExtendUint32WithBit(lower12Bits, 11) + ex.Get(reg)
	ex.operator.storeByte(src, address, memory)
}

/*CsrReadAndWrite atomically reads the value of CSR `csr` into register `dest` and writes the value of
register `reg` into CSR `csr`. However, the read does not occur AT ALL if `dest` == 0*/
func (ex *RiscVInstructionExecutor) CsrReadAndWrite(dest uint, reg uint, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	regVal := ex.Get(reg)

	if dest != 0 {
		csrVal := csrOperator.get(csr)
		ex.operator.andImmediate(dest, dest, 0)
		ex.operator.orImmediate(dest, dest, csrVal)
	}

	csrOperator.set(csr, regVal)
}

/*CsrReadAndSet atomically reads the value of CSR `csr` into register `dest`, performs a logical
OR between the value in `csr` and the value in `reg`, and the result is written back into the CSR `csr`
However, the write to `csr` will not happen AT ALL if `reg` == 0*/
func (ex *RiscVInstructionExecutor) CsrReadAndSet(dest uint, reg uint, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	regVal := ex.Get(reg)
	csrVal := csrOperator.get(csr)

	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, csrVal)

	if reg != 0 {
		csrOperator.set(csr, regVal|csrVal)
	}

}

/*CsrReadAndClear atomically reads the value of CSR `csr` into register `dest`, and uses the high bits in
`reg` to clear the corresponding bits in `csr`. However, the write to `csr` will not happen AT ALL if `reg` == 0*
*/
func (ex *RiscVInstructionExecutor) CsrReadAndClear(dest uint, reg uint, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	regVal := ex.Get(reg)
	csrVal := csrOperator.get(csr)

	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, csrVal)

	if reg != 0 {
		csrOperator.set(csr, (^regVal)&csrVal)
	}
}

/*CsrReadAndWriteImmediate atomically reads the value of CSR `csr` into register `dest` and writes zero-extended
lowest 5 bits of `immediate` into CSR `csr`. However, the read does not occur AT ALL if `dest` == 0
*/
func (ex *RiscVInstructionExecutor) CsrReadAndWriteImmediate(dest uint, immediate uint32, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	if dest != 0 {
		csrVal := csrOperator.get(csr)
		ex.operator.andImmediate(dest, dest, 0)
		ex.operator.orImmediate(dest, dest, csrVal)
	}

	csrOperator.set(csr, Utils.KeepBitsInInclusiveRange(immediate, 0, 4))

}

/*CsrReadAndSetImmediate atomically reads the value of CSR `csr` into register `dest`, performs a logical
OR between the value in `csr` and the zero-extended lowest 5 bits of `immediate`,
and the result is written back into the CSR `csr`
However, the write to `csr` will not happen AT ALL if the lowest 5 bits of `immediate` == 0*/
func (ex *RiscVInstructionExecutor) CsrReadAndSetImmediate(dest uint, immediate uint32, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	immVal := Utils.KeepBitsInInclusiveRange(immediate, 0, 4)
	csrVal := csrOperator.get(csr)

	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, csrVal)

	if immVal != 0 {
		csrOperator.set(csr, immVal|csrVal)
	}

}

/*CsrReadAndClearImmediate atomically reads the value of CSR `csr` into register `dest`, and uses the high bits in
the zero-extended lowest 5 bits of `immediate` to clear the corresponding bits in `csr`.
However, the write to `csr` will not happen AT ALL if the lowest 5 bits of `immediate` == 0*
*/
func (ex *RiscVInstructionExecutor) CsrReadAndClearImmediate(dest uint, immediate uint32, csr uint, csrOperator csrOperator) {
	defer ex.resetRegisterZero()

	immVal := Utils.KeepBitsInInclusiveRange(immediate, 0, 4)
	csrVal := csrOperator.get(csr)

	ex.operator.andImmediate(dest, dest, 0)
	ex.operator.orImmediate(dest, dest, csrVal)

	if immVal != 0 {
		csrOperator.set(csr, (^immVal)&csrVal)
	}
}

/*EnvCall is used to user-level programs to make a request to the supporting execution environment, which is usually an operating system.
According to the RiscV spec, how the arguments are passed is up to the execution environment, but will
typically be in a defined set of the registers.
*/
func (ex *RiscVInstructionExecutor) EnvCall(env executionEnvManager) {
	defer ex.resetRegisterZero()

	env.executeCall()
}

/*EnvBreak is used by user-level level programs to transfer control to a supporting debugging environment.
This seems like it is similar to a the EnvCall, in that arguments can be passed in the registers according to the
debugging environment's ABI. */
func (ex *RiscVInstructionExecutor) EnvBreak(env debugEnvManager) {
	defer ex.resetRegisterZero()

	env.debugBreak()
}
