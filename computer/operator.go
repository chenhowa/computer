package computer

type Operator struct {
	registers [16]uint32
	flags     uint16
}

type ReadMemory interface {
	Get(address uint16) uint32
}

type WriteMemory interface {
	Set(address uint16, value uint32)
}

func (c *Operator) load(reg uint, address uint16, m ReadMemory) {
	c.registers[reg] = m.Get(address)
}

func (c *Operator) store(reg uint, address uint16, m WriteMemory) {
	m.Set(address, c.registers[reg])
}

/*
	This function needs check for overflow.
*/
func (c *Operator) add(dest uint, reg1 uint, reg2 uint) {
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	var result = operand_1 + operand_2

	c.registers[dest] = result
}

/*
	This function needs check for overflow.
*/
func (c *Operator) sub(dest uint, reg1 uint, reg2 uint) {
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	var result = operand_1 - operand_2

	c.registers[dest] = result
}

func (c *Operator) bit_and(dest uint, reg1 uint, reg2 uint) {
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	var result = operand_1 & operand_2

	c.registers[dest] = result
}

func (c *Operator) bit_or(dest uint, reg1 uint, reg2 uint) {
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	var result = operand_1 | operand_2

	c.registers[dest] = result
}

/*
	This function should flag the msb
*/
func (c *Operator) left_shift(dest uint, reg uint) {
	var operand uint32 = c.registers[reg]
	//var msb = operand >> 31;
	var result = operand << 1

	c.registers[dest] = result
}

/*
	This function should flag the lsb
*/
func (c *Operator) right_shift(dest uint, reg uint, preserve_sign bool) {
	var operand uint32 = c.registers[reg]
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
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	var result = operand_1 * operand_2

	c.registers[dest] = result
}

/*
	This function should flag the underflow.
*/
func (c *Operator) divide(dest_dividend uint, dest_rem uint, reg1 uint, reg2 uint) {
	var operand_1 uint32 = c.registers[reg1]
	var operand_2 uint32 = c.registers[reg2]
	dividend := operand_1 / operand_2
	remainder := operand_1 % operand_2

	c.registers[dest_dividend] = dividend
	c.registers[dest_rem] = remainder
}
