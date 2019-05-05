package instructionParsing

import "fmt"

/*RiscVBinaryInstructionParser parses 32-bit binary instructions
into the appropriate fields
*/
type RiscVBinaryInstructionParser struct {
}

/*
	Will need to write a connector between this parser and the assembly executor
	that will use it.
*/
type OpCode uint

/*These constants represent the various opcodes of instructions
 */
const (
	ImmArith OpCode = iota
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

/*Parse will take a 32 bit instruction and parse its
fields into the relevant values BY INSTRUCTION TYPE,
and then return the results in a Struct.
*/
func (parser *RiscVBinaryInstructionParser) Parse(instruction uint32) RiscVBinaryParseResult {
	opcode := OpCode(((1 << 7) - 1) & instruction)
	var result RiscVBinaryParseResult

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
	result.OpCode = opcode

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

func parseAsI(instruction uint32) RiscVBinaryParseResult {
	var uintInstruction = uint(instruction)
	result := RiscVBinaryParseResult{InstructionType: I}
	result.FiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))
	result.Funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.FiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.TwelveBitImmediate = uint16(getBitsInInclusiveRange(uintInstruction, 20, 31))

	return result
}

func parseAsU(instruction uint32) RiscVBinaryParseResult {
	var uintInstruction = uint(instruction)
	result := RiscVBinaryParseResult{InstructionType: U}
	result.FiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))
	result.TwentyBitImmediate = uint32(getBitsInInclusiveRange(uintInstruction, 12, 31))

	return result
}

func parseAsR(instruction uint32) RiscVBinaryParseResult {
	result := RiscVBinaryParseResult{InstructionType: R}
	var uintInstruction = uint(instruction)
	result.FiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))
	result.Funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.FiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.FiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))
	result.Funct7 = uint8(getBitsInInclusiveRange(uintInstruction, 25, 31))

	return result
}

func parseAsJ(instruction uint32) RiscVBinaryParseResult {
	result := RiscVBinaryParseResult{InstructionType: J}
	var uintInstruction = uint(instruction)
	result.FiveBitDestination = uint8(getBitsInInclusiveRange(uintInstruction, 7, 11))

	bits1To10 := getBitsInInclusiveRange(uintInstruction, 21, 30)
	bit11 := getBitsInInclusiveRange(uintInstruction, 20, 20)
	bits12To19 := getBitsInInclusiveRange(uintInstruction, 12, 19)
	bit20 := getBitsInInclusiveRange(uintInstruction, 31, 31)

	result.TwentyBitImmediate = uint32(bits1To10 | (bit11 << 10) | (bits12To19 << 11) | (bit20 << 19))
	return result
}

func parseAsS(instruction uint32) RiscVBinaryParseResult {
	result := RiscVBinaryParseResult{InstructionType: S}
	var uintInstruction = uint(instruction)
	result.Funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.FiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.FiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))

	lowerFiveBits := getBitsInInclusiveRange(uintInstruction, 7, 11)
	upper7Bits := getBitsInInclusiveRange(uintInstruction, 25, 31)
	result.TwelveBitImmediate = uint16(lowerFiveBits | (upper7Bits << 5))

	return result
}

func parseAsB(instruction uint32) RiscVBinaryParseResult {
	result := RiscVBinaryParseResult{InstructionType: B}
	var uintInstruction = uint(instruction)
	result.Funct3 = uint8(getBitsInInclusiveRange(uintInstruction, 12, 14))
	result.FiveBitRegister1 = uint8(getBitsInInclusiveRange(uintInstruction, 15, 19))
	result.FiveBitRegister2 = uint8(getBitsInInclusiveRange(uintInstruction, 20, 24))

	bits1To4 := getBitsInInclusiveRange(uintInstruction, 8, 11)
	bits5To10 := getBitsInInclusiveRange(uintInstruction, 25, 30)
	bit11 := getBitsInInclusiveRange(uintInstruction, 7, 7)
	bit12 := getBitsInInclusiveRange(uintInstruction, 31, 31)
	result.TwelveBitImmediate = uint16(bits1To4 | (bits5To10 << 4) | (bit11 << 10) | (bit12 << 11))

	return result
}

/*InstructionType represents valid Instruction types for 32I RiscV, where each Instruction Type
defines a format for representing the operands of the 32-bit instruction. The only way to
know an instruction type is to check the OpCode of the instruction -- each OpCode is mapped
to just one available Instruction Type.
*/
type InstructionType uint

/*These constants represent the various available instruction types available in the RiscV instruction set (32 bits)
 */
const (
	R InstructionType = iota // the below instruction types are the base instruction types
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

/*RiscVBinaryParseResult exposes the parsed binary
instruction. It has multiple fields; not all are valid
for a given Instruction Type. The Valid fields for a
given instruction type follow the RiscV spec.
*/
type RiscVBinaryParseResult struct {
	InstructionType    InstructionType
	OpCode             OpCode
	FiveBitRegister1   uint8
	FiveBitRegister2   uint8
	FiveBitDestination uint8
	TwelveBitImmediate uint16
	Funct3             uint8
	Funct7             uint8
	TwentyBitImmediate uint32
}

// the purpose of this function is to check that all values that populate
// the result fit the constraints. If they do not, it is a runtime error
// Unfortunately we cannot catch these with the OOB type system.
func (result RiscVBinaryParseResult) errorIfInvalid() {
	if xIsGreaterThanYBits(uint(result.FiveBitRegister1), 5) {
		panic(fmt.Sprintf("invalid binary parse result FiveBitRegister1 %d", result.FiveBitRegister1))
	} else if xIsGreaterThanYBits(uint(result.FiveBitRegister2), 5) {
		panic(fmt.Sprintf("invalid binary parse result FiveBitRegister2 %d", result.FiveBitRegister2))
	} else if xIsGreaterThanYBits(uint(result.FiveBitDestination), 5) {
		panic(fmt.Sprintf("invalid binary parse result FiveBitDestination %d", result.FiveBitDestination))
	} else if xIsGreaterThanYBits(uint(result.TwelveBitImmediate), 12) {
		panic(fmt.Sprintf("invalid binary parse result TwelveBitImmediate %d", result.TwelveBitImmediate))
	} else if xIsGreaterThanYBits(uint(result.Funct3), 3) {
		panic(fmt.Sprintf("invalid binary parse result Funct3 %d", result.Funct3))
	} else if xIsGreaterThanYBits(uint(result.Funct7), 7) {
		panic(fmt.Sprintf("invalid binary parse result Funct7 %d", result.Funct7))
	} else if xIsGreaterThanYBits(uint(result.TwentyBitImmediate), 20) {
		panic(fmt.Sprintf("invalid binary parse result TwentyBitImmediate %d", result.TwentyBitImmediate))
	}
}

func xIsGreaterThanYBits(x uint, y uint) bool {
	return x >= (1 << y)
}
