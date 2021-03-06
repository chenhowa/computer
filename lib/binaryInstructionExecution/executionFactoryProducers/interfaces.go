package executionFactoryProducers

/*The RiscVExecutor describes a receiver that is capable of
executing all 32I RiscV instructions, using the fields encoded
in each type of instruction
*/
type RiscVExecutor interface {
	addImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32)
	andImmmediate(dest uint, reg uint, immediate uint32)
	orImmediate(dest uint, reg uint, immediate uint32)
	xorImmediate(dest uint, reg uint, immediate uint32)
	shiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32)
	shiftRight(dest uint, reg uint, immediate uint32)
	loadUpperImmediate(dest uint, immediate uint32)
	addUpperImmediateToPC(dest uint, immediate uint32)
	add(dest uint, reg1 uint, reg2 uint)
	sub(dest uint, reg1 uint, reg2 uint)
	setLessThan(dest uint, reg1 uint, reg2 uint)
	setLessThanUnsigned(dest uint, reg1 uint, reg2 uint)
	and(dest uint, reg1 uint, reg2 uint)
	or(dest uint, reg1 uint, reg2 uint)
	xor(dest uint, reg1 uint, reg2 uint)
	shiftLeftLogical(dest uint, reg uint, shiftreg uint)
	shiftRightLogical(dest uint, reg uint, shiftreg uint)
	shiftRightArithmetic(dest uint, reg uint, shiftreg uint)
	branchEqual(src1 uint, src2 uint, offset uint32)
	branchNotEqual(src1 uint, src2 uint, offset uint32)
	branchLessThan(src1 uint, src2 uint, offset uint32)
	branchLessThanUnsigned(src1 uint, src2 uint, offset uint32)
	branchGreaterThanOrEqual(src1 uint, src2 uint, offset uint32)
	branchGreaterThanOrEqualUnsigned(src1 uint, src2 uint, offset uint32)
	jumpAndLink(dest uint, pcOffset uint32)
	jumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32)
	loadWord(dest uint, reg uint, offset uint32)
	loadHalfWord(dest uint, reg uint, offset uint32)
	loadHalfWordUnsigned(dest uint, reg uint, offset uint32)
	loadByte(dest uint, reg uint, offset uint32)
	loadByteUnsigned(dest uint, reg uint, offset uint32)
	storeWord(reg1 uint, reg2 uint, offset uint32)
	storeHalfWord(reg1 uint, reg2 uint, offset uint32)
	storeByte(reg1 uint, reg2 uint, offset uint32)
	csrReadAndWrite(dest uint, reg uint, immediate uint32)
	csrReadAndSet(dest uint, reg uint, immediate uint32)
	csrReadAndClear(dest uint, reg uint, immediate uint32)
	csrReadAndWriteImmediate(dest uint, reg uint, immediate uint32)
	csrReadAndSetImmediate(dest uint, reg uint, immediate uint32)
	csrReadAndClearImmediate(dest uint, reg uint, immediate uint32)
	private(dest uint, reg uint, immediate uint32)
}
