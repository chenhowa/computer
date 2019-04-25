package computer

import "fmt"

type riscVBinaryInstructionParser struct {
}

func (parser *riscVBinaryInstructionParser) parse(instruction uint32) riscVBinaryParseResult {
	opcode := 127 & instruction
	var instructionType instructionType
	var result riscVBinaryParseResult

	if opcode < 3 { // immediate arithmetic
		result = parseAsI(instruction)
	} else if opcode < 4 { // register-register arithmetic
		result = parseAsR(instruction)
	} else if opcode < 5 { // jump and link
		result = parseAsJ(instruction)
	} else if opcode < 6 { // jump and link register
		result = parseAsI(instruction)
	} else if opcode < 7 { // branch
		result = parseAsB(instruction)
	} else if opcode < 8 { // loads
		result = parseAsI(instruction)
	} else if opcode < 9 { // stores
		result = parseAsS(instruction)
	} else if opcode < 10 { // system (csr)
		result = parseAsI(instruction)
	} else {
		panic(fmt.Sprintf("unrecognized opcode %d", opcode))
	}

	return riscVBinaryParseResult{}
}

func parseAsI(instruction uint32) riscVBinaryParseResult {

}

func parseAsR(instruction uint32) riscVBinaryParseResult {

}

func parseAsJ(instruction uint32) riscVBinaryParseResult {

}

func parseAsS(instruction uint32) riscVBinaryParseResult {

}

func parseAsB(instruction uint32) riscVBinaryParseResult {

}

type instructionType uint

const (
	// the below instruction types are the base instruction types
	r instructionType = iota + 1
	i
	s
	u
	b

	// the below constants are extension instruction types. see the online specification for more.
	m
	a
	f
	d
	q
	l
	c
	j
	t
	p
	v
	n
)

type riscVBinaryParseResult struct {
	instructionType    instructionType
	fiveBitRegister1   uint8
	fiveBitRegister2   uint8
	fiveBitDestination uint8
	twelveBitImmediate uint16
	funct3             uint8
	funct7             uint8
	twentyBitImmediate uint32
}

// the purpose of this function is to check that all values that populate
// the result fit the constraints. If they do not, it is a runtime error
// Unfortunately we cannot catch these with the OOB type system.
func (result riscVBinaryParseResult) errorIfInvalid() {
	if xIsGreaterThanYBits(uint(result.fiveBitRegister1), 5) {
		panic(fmt.Sprintf("invalid binary parse result fiveBitRegister1 %d", result.fiveBitRegister1))
	} else if xIsGreaterThanYBits(uint(result.fiveBitRegister2), 5) {
		panic(fmt.Sprintf("invalid binary parse result fiveBitRegister2 %d", result.fiveBitRegister2))
	} else if xIsGreaterThanYBits(uint(result.fiveBitDestination), 5) {
		panic(fmt.Sprintf("invalid binary parse result fiveBitDestination %d", result.fiveBitDestination))
	} else if xIsGreaterThanYBits(uint(result.twelveBitImmediate), 12) {
		panic(fmt.Sprintf("invalid binary parse result twelveBitImmediate %d", result.twelveBitImmediate))
	} else if xIsGreaterThanYBits(uint(result.funct3), 3) {
		panic(fmt.Sprintf("invalid binary parse result funct3 %d", result.funct3))
	} else if xIsGreaterThanYBits(uint(result.funct7), 7) {
		panic(fmt.Sprintf("invalid binary parse result funct7 %d", result.funct7))
	} else if xIsGreaterThanYBits(uint(result.twentyBitImmediate), 20) {
		panic(fmt.Sprintf("invalid binary parse result twentyBitImmediate %d", result.twentyBitImmediate))
	}
}

func xIsGreaterThanYBits(x uint, y uint) bool {
	return x >= (1 >> y)
}

/*
	Will need to write a connector between this parser and the assembly executor
	that will use it.
*/
