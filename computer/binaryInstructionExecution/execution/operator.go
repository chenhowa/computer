package execution

import (
	"fmt"
	"math"
)

/*Operator represents operations a set of 16(??) registers.

 */
type Operator struct {
	registers [16]uint32
	flags     uint16
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

func (c *Operator) get(reg uint) uint32 {
	return c.registers[reg]
}

func (c *Operator) load(reg uint, address uint16, m ReadMemory) {
	c.registers[reg] = m.Get(address)
}

func (c *Operator) store(reg uint, address uint16, m WriteMemory) {
	m.Set(address, c.registers[reg], 32)
}

func (c *Operator) add_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] + immediate
}

func (c *Operator) bit_and_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] & immediate
}

func (c *Operator) bit_or_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] | immediate
}

func (c *Operator) bit_xor_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] ^ immediate
}

func (c *Operator) bit_not_immediate(dest uint, immediate uint32) {
	c.registers[dest] = ^immediate
}

func (c *Operator) left_shift_immediate(dest uint, reg uint, immediate uint32) {
	c.registers[dest] = c.registers[reg] << immediate
}

func (c *Operator) right_shift_immediate(dest uint, reg uint, immediate uint32, preserveSign bool) {
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
func (c *Operator) add(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 + operand2

	c.registers[dest] = result
}

/*
	This function needs check for overflow.
*/
func (c *Operator) sub(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 - operand2

	c.registers[dest] = result
}

func (c *Operator) bit_and(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 & operand2

	c.registers[dest] = result
}

func (c *Operator) bit_or(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 | operand2

	c.registers[dest] = result
}

func (c *Operator) bit_xor(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 ^ operand2

	c.registers[dest] = result
}

func (c *Operator) bit_not(dest uint, reg1 uint) {
	var operand1 = c.registers[reg1]
	var result = ^operand1

	c.registers[dest] = result
}

/*
	This function should flag the msb
*/
func (c *Operator) left_shift(dest uint, reg uint) {
	var operand = c.registers[reg]
	//var msb = operand >> 31;
	var result = operand << 1

	c.registers[dest] = result
}

/*
	This function should flag the lsb
*/
func (c *Operator) right_shift(dest uint, reg uint, preserve_sign bool) {
	var operand = c.registers[reg]
	//var lsb = operand & 1;
	var result uint32

	// Go automatically arithmetic shifts if the result is signed.
	if preserve_sign {
		result = uint32(int32(operand) >> 1)
	} else {
		result = operand >> 1
	}

	c.registers[dest] = result
}

/*
	This function should flag the overflow.
*/
func (c *Operator) multiply(dest uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	var result = operand1 * operand2

	c.registers[dest] = result
}

/*
	This function should flag the underflow.
*/
func (c *Operator) divide(dest_dividend uint, dest_rem uint, reg1 uint, reg2 uint) {
	var operand1 = c.registers[reg1]
	var operand2 = c.registers[reg2]
	dividend := operand1 / operand2
	remainder := operand1 % operand2

	c.registers[dest_dividend] = dividend
	c.registers[dest_rem] = remainder
}

func keepBitsInInclusiveRange(num uint32, start uint, end uint) uint32 {
	if start < end {
		panic(fmt.Sprintf("keepBitsInRange: start > end, %d > %d", start, end))
	}

	majorMask := uint32(1<<(end+1)) - 1
	minorMask := uint32(math.MaxUint32) << (start + 1)

	return (num & majorMask & minorMask)
}

func (c *Operator) load_word(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = value
}

func (c *Operator) load_halfword(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = signExtendUint32WithBit(keepBitsInInclusiveRange(value, 0, 15), 15)
}

func (c *Operator) load_halfword_unsigned(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = keepBitsInInclusiveRange(value, 0, 15)
}

func (c *Operator) load_byte(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = signExtendUint32WithBit(keepBitsInInclusiveRange(value, 0, 7), 7)
}

func (c *Operator) load_byte_unsigned(dest uint, address uint16, memory ReadMemory) {
	value := memory.Get(address)
	c.registers[dest] = keepBitsInInclusiveRange(value, 0, 7)
}

func (c *Operator) store_word(src uint, address uint16, memory WriteMemory) {
	value := c.get(src)
	memory.Set(address, value, 32)
}

func (c *Operator) store_halfword(src uint, address uint16, memory WriteMemory) {
	value := c.get(src)
	memory.Set(address, value, 16)
}

func (c *Operator) store_byte(src uint, address uint16, memory WriteMemory) {
	value := c.get(src)
	memory.Set(address, value, 8)
}
