package execution

import (
	"math"

	Operator "github.com/chenhowa/computer/lib/binaryInstructionExecution/execution/operators"
)

/*adaptedOperator is a an adapter for the imported Operator struct to help fit
the interface that the RiscVInstructionExecutor requires.
*/
type adaptedOperator struct {
	operator *Operator.Operator
}

func makeAdaptedOperator(regs [32]uint32) adaptedOperator {
	op := Operator.MakeOperator(regs, 0)

	adapted := adaptedOperator{
		operator: &op,
	}

	return adapted
}

/*adaptedMemory is an adapter for structs of types `instructionReadMemory`, `instructionWriteMemory`,
and `instructionReadWriteMemory`, to help the struct fit the constraints of the imported Operator struct.
*/
type adaptedMemory struct {
	rMemory instructionReadMemory
	wMemory instructionWriteMemory
}

func (m *adaptedMemory) Get(address uint16) uint32 {
	return m.rMemory.Get(uint32(address))
}

func (m *adaptedMemory) Set(address uint16, value uint32, bitsToSet uint) {
	m.wMemory.Set(uint32(address), value, bitsToSet)
}

/*panicIfOutsideMemory will panic if the uint32 `address` cannot be casted to
uint16 without losing bits. Since the Operator only supports 16-bit addresses.
*/
func panicIfOutsideMemory(address uint32) {
	if address > uint32(math.MaxUint16) {
		panic("address outside uint16 actual memory space")
	}
}

func (op *adaptedOperator) loadWord(dest uint, address uint32, memory instructionReadMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		rMemory: memory,
	}
	op.operator.Load_word(dest, uint16(address), &m)

}

func (op *adaptedOperator) loadHalfWord(dest uint, address uint32, memory instructionReadMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		rMemory: memory,
	}
	op.operator.Load_halfword(dest, uint16(address), &m)
}

func (op *adaptedOperator) loadHalfWordUnsigned(dest uint, address uint32, memory instructionReadMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		rMemory: memory,
	}
	op.operator.Load_halfword_unsigned(dest, uint16(address), &m)
}

func (op *adaptedOperator) loadByte(dest uint, address uint32, memory instructionReadMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		rMemory: memory,
	}
	op.operator.Load_byte(dest, uint16(address), &m)
}

func (op *adaptedOperator) loadByteUnsigned(dest uint, address uint32, memory instructionReadMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		rMemory: memory,
	}
	op.operator.Load_byte_unsigned(dest, uint16(address), &m)
}

func (op *adaptedOperator) storeWord(src uint, address uint32, memory instructionWriteMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		wMemory: memory,
	}
	op.operator.Store_word(src, uint16(address), &m)
}

func (op *adaptedOperator) storeHalfWord(src uint, address uint32, memory instructionWriteMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		wMemory: memory,
	}
	op.operator.Store_word(src, uint16(address), &m)
}

func (op *adaptedOperator) storeByte(src uint, address uint32, memory instructionWriteMemory) {
	panicIfOutsideMemory(address)
	m := adaptedMemory{
		wMemory: memory,
	}
	op.operator.Store_word(src, uint16(address), &m)
}

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
