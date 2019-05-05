package execution

import (
	Operator "github.com/chenhowa/os/computer/binaryInstructionExecution/execution/operators"
)

/*adaptedOperator is a an adapter for the imported Operator struct to help fit
the interface that the RiscVInstructionExecutor requires.
*/
type adaptedOperator struct {
	operator Operator.Operator
}

func (op *adaptedOperator) loadWord(dest uint, address uint32, memory instructionReadMemory)     {}
func (op *adaptedOperator) loadHalfWord(dest uint, address uint32, memory instructionReadMemory) {}
func (op *adaptedOperator) loadHalfWordUnsigned(dest uint, address uint32, memory instructionReadMemory) {
}
func (op *adaptedOperator) loadByte(dest uint, address uint32, memory instructionReadMemory)         {}
func (op *adaptedOperator) loadByteUnsigned(dest uint, address uint32, memory instructionReadMemory) {}
func (op *adaptedOperator) storeWord(src uint, address uint32, memory instructionWriteMemory)        {}
func (op *adaptedOperator) storeHalfWord(src uint, address uint32, memory instructionWriteMemory)    {}
func (op *adaptedOperator) storeByte(src uint, address uint32, memory instructionWriteMemory)        {}
func (op *adaptedOperator) add(dest uint, reg1 uint, reg2 uint) {
	op.operator.Add(dest, reg1, reg2)
}
func (op *adaptedOperator) addImmediate(dest uint, reg uint, immediate uint32) {
	op.operator.Add_immediate(dest, reg, immediate)
}
func (op *adaptedOperator) sub(dest uint, reg1 uint, reg2 uint) {
	op.operator.Sub(dest, reg1, reg2)
}
func (op *adaptedOperator) and(dest uint, reg1 uint, reg2 uint) {
	op.operator.Bit_and(dest, reg1, reg2)
}
func (op *adaptedOperator) andImmediate(dest uint, reg uint, immediate uint32) {
	op.operator.Bit_and_immediate(dest, reg, immediate)
}
func (op *adaptedOperator) or(dest uint, reg1 uint, reg2 uint) {
	op.operator.Bit_or(dest, reg1, reg2)
}
func (op *adaptedOperator) orImmediate(dest uint, reg uint, immediate uint32) {
	op.operator.Bit_or_immediate(dest, reg, immediate)
}

func (op *adaptedOperator) xor(dest uint, reg1 uint, reg2 uint) {
	op.operator.Bit_xor(dest, reg1, reg2)
}

func (op *adaptedOperator) xorImmediate(dest uint, reg uint, immediate uint32) {
	op.operator.Bit_xor_immediate(dest, reg, immediate)
}

func (op *adaptedOperator) leftShiftImmediate(dest uint, reg uint, immediate uint32) {
	op.operator.Left_shift_immediate(dest, reg, immediate)
}

func (op *adaptedOperator) rightShiftImmediate(dest uint, reg uint, immediate uint32, preserveSign bool) {
	op.operator.Right_shift_immediate(dest, reg, immediate, preserveSign)
}

func (op *adaptedOperator) multiply(dest uint, reg1 uint, reg2 uint) {
	op.operator.Multiply(dest, reg1, reg2)
}

func (op *adaptedOperator) divide(destDividend uint, destRem uint, reg1 uint, reg2 uint) {
	op.operator.Divide(destDividend, destRem, reg1, reg2)
}

func (op *adaptedOperator) get(reg uint) uint32 {
	return op.operator.Get(reg)
}
