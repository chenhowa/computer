package operators

import (
	Utils "github.com/chenhowa/os/computer/binaryInstructionExecution/bitUtils"
)

/*Operator represents operations a set of 16(??) registers.

 */
type Operator struct {
	registers [16]uint32
	flags     uint16
}

/*MakeOperator is a construction function for the Operator struct
 */
func MakeOperator(registers [16]uint32, flags uint16) Operator {
	op := Operator{
		registers: registers,
		flags:     flags,
	}

	return op
}

/*ReadMemory is an interface that reads a uint32 data
value from a memory address (possibly virtualized)
*/
type ReadMemory interface {
	Get(address uint16) uint32
}

/*WriteMemory is an interface that writes a uint32 data
value to a memory address (possibly virtualized)
*/
type WriteMemory interface {
	Set(address uint16, value uint32, bitsToSet uint)
}

type ReadWriteMemory interface {
	ReadMemory
	WriteMemory
}

func (c *Operator) Get(reg uint) uint32 {
	return c.registers[reg]
}

func (c *Operator) Load(reg uint, address uint16, m ReadMemory) {
	c.registers[reg] = m.Get(address)
}

func (c *Operator) Store(reg uint, address uint16, m WriteMemory) {
	m.Set(address, c.registers[reg], 32)
}

func (c *Operator) Add_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] + immediate
}

func (c *Operator) Bit_and_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] & immediate
}

func (c *Operator) Bit_or_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] | immediate
}

func (c *Operator) Bit_xor_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] ^ immediate
}

func (c *Operator) Bit_not_immediate(dest uint, immediate uint32) {
	c.registers[dest] = ^immediate
}

func (c *Operator) Left_shift_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] << immediate
}

func (c *Operator) Right_shift_immediate(dest uint, reg uint, immediate uint32, preserveSign bool) {
	operand := c.registers[reg]
	var result uint32

	// Go automatically arithmetic shifts if the result is signed.
	if preserveSign {
		result = uint32(int32(operand) >> immediate)
	} else {
		result = operand >> immediate
	}

	c.registers[dest] = result
}

/*
	This function needs check for overflow.
*/
func (c *Operator) Add(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 + operand2

	c.registers[dest] = result
}

/*
	This function needs check for overflow.
*/
func (c *Operator) Sub(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 - operand2

	c.registers[dest] = result
}

func (c *Operator) Bit_and(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 & operand2

	c.registers[dest] = result
}

func (c *Operator) Bit_or(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 | operand2

	c.registers[dest] = result
}

func (c *Operator) Bit_xor(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 ^ operand2

	c.registers[dest] = result
}

func (c *Operator) Bit_not(dest uint, reg1 uint) {
	var operand1 = c.registers[reg1]
	var result = ^operand1

	c.registers[dest] = result
}

/*
	This function should flag the overflow.
*/
func (c *Operator) Multiply(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 * operand2

	c.registers[dest] = result
}

/*
	This function should flag the underflow.
*/
func (c *Operator) Divide(dest_dividend uint, dest_rem uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	dividend := operand1 / operand2
	remainder := operand1 % operand2

	c.registers[dest_dividend] = dividend
	c.registers[dest_rem] = remainder
}

func (c *Operator) Load_word(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = value
}

func (c *Operator) Load_halfword(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = Utils.SignExtendUint32WithBit(Utils.KeepBitsInInclusiveRange(value, 0, 15), 15)
}

func (c *Operator) Load_halfword_unsigned(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = Utils.KeepBitsInInclusiveRange(value, 0, 15)
}

func (c *Operator) Load_byte(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = Utils.SignExtendUint32WithBit(Utils.KeepBitsInInclusiveRange(value, 0, 7), 7)
}

func (c *Operator) Load_byte_unsigned(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = Utils.KeepBitsInInclusiveRange(value, 0, 7)
}

func (c *Operator) Store_word(src uint, address uint16, memory WriteMemory) {
	value := c.Get(src)
	memory.Set(address, value, 32)
}

func (c *Operator) Store_halfword(src uint, address uint16, memory WriteMemory) {
	value := c.Get(src)
	memory.Set(address, value, 16)
}

func (c *Operator) Store_byte(src uint, address uint16, memory WriteMemory) {
	value := c.Get(src)
	memory.Set(address, value, 8)
}
