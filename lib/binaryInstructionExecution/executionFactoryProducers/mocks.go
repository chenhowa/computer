package executionFactoryProducers

import "github.com/stretchr/testify/mock"

type RiscVExecutorMock struct {
	mock.Mock
}

func (em *RiscVExecutorMock) addImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) setLessThanImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) setLessThanImmediateUnsigned(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) andImmmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) orImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) xorImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) shiftLeftLogicalImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) shiftRight(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) loadUpperImmediate(dest uint, immediate uint32) {
	em.Called(dest, immediate)
}
func (em *RiscVExecutorMock) addUpperImmediateToPC(dest uint, immediate uint32) {
	em.Called(dest, immediate)
}
func (em *RiscVExecutorMock) add(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)
}
func (em *RiscVExecutorMock) sub(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)
}
func (em *RiscVExecutorMock) setLessThan(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)
}
func (em *RiscVExecutorMock) setLessThanUnsigned(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)

}
func (em *RiscVExecutorMock) and(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)

}
func (em *RiscVExecutorMock) or(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)

}
func (em *RiscVExecutorMock) xor(dest uint, reg1 uint, reg2 uint) {
	em.Called(dest, reg1, reg2)

}
func (em *RiscVExecutorMock) shiftLeftLogical(dest uint, reg uint, shiftreg uint) {
	em.Called(dest, reg, shiftreg)
}
func (em *RiscVExecutorMock) shiftRightLogical(dest uint, reg uint, shiftreg uint) {
	em.Called(dest, reg, shiftreg)

}
func (em *RiscVExecutorMock) shiftRightArithmetic(dest uint, reg uint, shiftreg uint) {
	em.Called(dest, reg, shiftreg)

}
func (em *RiscVExecutorMock) branchEqual(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)
}
func (em *RiscVExecutorMock) branchNotEqual(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)

}
func (em *RiscVExecutorMock) branchLessThan(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)

}
func (em *RiscVExecutorMock) branchLessThanUnsigned(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)

}
func (em *RiscVExecutorMock) branchGreaterThanOrEqual(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)

}
func (em *RiscVExecutorMock) branchGreaterThanOrEqualUnsigned(src1 uint, src2 uint, offset uint32) {
	em.Called(src1, src2, offset)

}
func (em *RiscVExecutorMock) jumpAndLink(dest uint, pcOffset uint32) {
	em.Called(dest, pcOffset)
}
func (em *RiscVExecutorMock) jumpAndLinkRegister(dest uint, basereg uint, pcOffset uint32) {
	em.Called(dest, basereg, pcOffset)
}
func (em *RiscVExecutorMock) loadWord(dest uint, reg uint, offset uint32) {
	em.Called(dest, reg, offset)
}
func (em *RiscVExecutorMock) loadHalfWord(dest uint, reg uint, offset uint32) {
	em.Called(dest, reg, offset)

}
func (em *RiscVExecutorMock) loadHalfWordUnsigned(dest uint, reg uint, offset uint32) {
	em.Called(dest, reg, offset)

}
func (em *RiscVExecutorMock) loadByte(dest uint, reg uint, offset uint32) {
	em.Called(dest, reg, offset)

}
func (em *RiscVExecutorMock) loadByteUnsigned(dest uint, reg uint, offset uint32) {
	em.Called(dest, reg, offset)

}
func (em *RiscVExecutorMock) storeWord(reg1 uint, reg2 uint, offset uint32) {
	em.Called(reg1, reg2, offset)
}
func (em *RiscVExecutorMock) storeHalfWord(reg1 uint, reg2 uint, offset uint32) {
	em.Called(reg1, reg2, offset)

}
func (em *RiscVExecutorMock) storeByte(reg1 uint, reg2 uint, offset uint32) {
	em.Called(reg1, reg2, offset)

}
func (em *RiscVExecutorMock) csrReadAndWrite(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)
}
func (em *RiscVExecutorMock) csrReadAndSet(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)

}
func (em *RiscVExecutorMock) csrReadAndClear(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)

}
func (em *RiscVExecutorMock) csrReadAndWriteImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)

}
func (em *RiscVExecutorMock) csrReadAndSetImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)

}
func (em *RiscVExecutorMock) csrReadAndClearImmediate(dest uint, reg uint, immediate uint32) {
	em.Called(dest, reg, immediate)

}
func (em *RiscVExecutorMock) private(dest uint, reg uint, immediate uint32) {
	em.private(dest, reg, immediate)
}
