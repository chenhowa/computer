package binaryInstructionExecution

import "fmt"

type riscVBinaryInstructionParser struct {
}

/*
	Will need to write a connector between this parser and the assembly executor
	that will use it.
*/
type opCode uint

/*These constants represent the various opcodes of instructions
 */
const (
	ImmArith opCode = iota + 1
	LUI
	AUIPC
	RegArith
	JAL
	JALR
	Branch
	Load
	Store
	System
)

func (parser *riscVBinaryInstructionParser) parse(instruction uint32) riscVBinaryParseResult {
	opcode := opCode(((1 << 7) - 1) & instruction)
	var result riscVBinaryParseResult

	if opcode < ImmArith+1 { // immediate arithmetic
		result = parseAsI(instruction)
	} else if opcode < AUIPC+1 {
		result = parseAsU(instruction)
	} else if opcode < RegArith+1 { // register-register arithmetic
		result = parseAsR(instruction)
	} else if opcode < JAL+1 { // jump and link
		result = parseAsJ(instruction)
	} else if opcode < JALR+1 { // jump and link register
		result = parseAsI(instruction)
	} else if opcode < Branch+1 { // branch
		result = parseAsB(instruction)
	} else if opcode < Load+1 { // loads
		result = parseAsI(instruction)
	} else if opcode < Store+1 { // stores
		result = parseAsS(instruction)
	} else if opcode < System+1 { // system (csr)
		result = parseAsI(instruction)
	} else {
		panic(fmt.Sprintf("unrecognized opcode %d", opcode))
	}
	result.errorIfInvalid()

	// Tag with the opcode, which was already calculated
	result.opCode = opcode

	return result
}

/*getBitsInInclusiveRange takes a number and extracts the value of a certain set of bits,
  from bit `start` to bit `end`. `start` is assumed to be less than `end`.
  The 0th bit is assumed to be the Least-Significant-Bit (LSB) of `number`

  Examples: when `number` = 0b000011100, `start` = 1, `end` = 4, getBitsInInclusiveRange(number, start, end) = 0b000001110 (as uint)
			when `number` = 0b000011100, `start` = 0, `end` = 2, getBitsInInclusiveRange(number, start, end) = 0b000000100 (as uint)
*/
func getBitsInInclusiveRange(number uint, start uint, end uint) uint {
	if start > end {
		panic(fmt.Sprintf("getBitsInInclusiveRange: start > end, start was %d, end was %d", start, end))
	}

	shifted := number >> start
	var masklength = end - start
	var mask uint = (1 << ((masklength) + 1)) - 1
	return shifted & mask
}

func parseAsI(instruction uint32) riscVBinaryParseResult {
	var uintInstruction = uint(instruction)
	result := riscVBinaryParseResult{instructionType: I}
	result.fiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))
	result.funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.fiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.twelveBitImmediate = uint16(getBitsInInclusiveRange(uintInstruction, 20, 31))

	return result
}

func parseAsR(instruction uint32) riscVBinaryParseResult {
	result := riscVBinaryParseResult{instructionType: R}
	var uintInstruction = uint(instruction)
	result.fiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))
	result.funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.fiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.fiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))
	result.funct7 = uint8(getBitsInInclusiveRange(uintInstruction, 25, 31))

	return result
}

func parseAsJ(instruction uint32) riscVBinaryParseResult {
	result := riscVBinaryParseResult{instructionType: J}
	var uintInstruction = uint(instruction)
	result.fiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))

	bits1To10 := getBitsInInclusiveRange(uintInstruction, 21, 30)
	bit11 := getBitsInInclusiveRange(uintInstruction, 20, 20)
	bits12To19 := getBitsInInclusiveRange(uintInstruction, 12, 19)
	bit20 := getBitsInInclusiveRange(uintInstruction, 31, 31)

	result.twentyBitImmediate = uint32(bits1To10 | (bit11 << 10) | (bits12To19 << 11) | (bit20 << 19))
	return result
}

func parseAsS(instruction uint32) riscVBinaryParseResult {
	result := riscVBinaryParseResult{instructionType: S}
	var uintInstruction = uint(instruction)
	result.funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.fiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.fiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))

	lowerFiveBits := getBitsInInclusiveRange(uintInstruction, 0, 4)
	upper7Bits := getBitsInInclusiveRange(uintInstruction, 25, 31)
	result.twelveBitImmediate = uint16(lowerFiveBits | (upper7Bits << 5))

	return result
}

func parseAsB(instruction uint32) riscVBinaryParseResult {
	result := riscVBinaryParseResult{instructionType: B}
	var uintInstruction = uint(instruction)
	result.funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.fiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.fiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))

	bits1To4 := getBitsInInclusiveRange(uintInstruction, 8, 11)
	bits5To10 := getBitsInInclusiveRange(uintInstruction, 25, 30)
	bit11 := getBitsInInclusiveRange(uintInstruction, 11, 11)
	bit12 := getBitsInInclusiveRange(uintInstruction, 31, 31)
	result.twelveBitImmediate = uint16(bits1To4 | (bits5To10 << 4) | (bit11 << 10) | (bit12 << 11))

	return result
}

type instructionType uint

/*These constants represent the various available instruction types available in the RiscV instruction set (32 bits)
 */
const (
	R instructionType = iota + 1 // the below instruction types are the base instruction types
	I
	S
	U
	B

	M // the below constants are extension instruction types. see the online specification for more.
	A
	F
	D
	Q
	L
	C
	J
	T
	P
	V
	N
)

type riscVBinaryParseResult struct {
	instructionType    instructionType
	opCode             opCode
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
