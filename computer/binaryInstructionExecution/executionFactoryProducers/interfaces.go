package executionFactoryProducers

type RiscVExecutor interface {
	addImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediate(dest uint, reg uint, immediate uint32)
	setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32)
	andImmmediate(dest uint, reg uint, immediate uint32)
	orImmediate(dest uint, reg uint, immediate uint32)
	xorImmediate(dest uint, reg uint, immediate uint32)
	shiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32)
	shiftRightLogicalImmediate(dest uint, reg uint, immediate uint32)
	shiftRightArithmeticImmediate(dest uint, reg uint, immediate uint32)
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
	branchEqual()
	branchNotEqual()
	branchLessThan()
	branchLessThanUnsigned()
	branchGreaterThanOrEqual()
	branchGreaterThanOrEqualUnsigned()
	jumpAndLink(dest uint, pcOffset uint32)
	jumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32)
	loadWord(dest uint, reg uint, offset uint32)
	loadHalfWord(dest uint, reg uint, offset uint32)
	loadHalfWordUnsigned(dest uint, reg uint, offset uint32)
	loadByte(dest uint, reg uint, offset uint32)
	loadByteUnsigned(dest uint, reg uint, offset uint32)
	storeWord(src uint, reg uint, offset uint32)
	storeHalfWord(src uint, reg uint, offset uint32)
	storeByte(src uint, reg uint, offset uint32)
	csrReadAndWrite()
	csrReadAndSet()
	csrReadAndClear()
	csrReadAndWriteImmediate()
	csrReadAndSetImmediate()
	csrReadAndClearImmediate()
	envCall()
	envBreak()
}
