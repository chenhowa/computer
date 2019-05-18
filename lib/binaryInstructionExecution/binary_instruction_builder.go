package binaryInstructionExecution

import (
	Utils "github.com/chenhowa/computer/lib/binaryInstructionExecution/bitUtils"
	Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"
)

/*BuildInstructionI builds a 32-bit I Format instruction out of the arguments.
Uses lowest bits of arguments as follows:
	- 7 bits of `opcode`
	- 5 bits of `rd`
	- 3 bits `funct3`
	- 5 bits of `rs1`
	- 12 bits of `immediate`
*/
func BuildInstructionI(opcode uint, rd uint, funct3 uint, rs1 uint, immediate uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, rd)
	builder.AddNextXBits(3, funct3)
	builder.AddNextXBits(5, rs1)
	builder.AddNextXBits(12, immediate)
	instruction := builder.Build()
	return uint32(instruction)
}

/*BuildInstructionU builds a 32-bit U Format instruction out of the arguments.
Uses lowest bits of arguments as follows:
	- 7 bits of `opcode`
	- 5 bits of `rd`
	- 20 bits of `immediate`*/
func BuildInstructionU(opcode uint, rd uint, immediate uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, rd)
	builder.AddNextXBits(20, immediate)
	return uint32(builder.Build())
}

/*BuildInstructionR builds a 32-bit R Format instruction out of the arguments.
Uses lowest bits of arguments as follows:
	- 7 bits of `opcode`
	- 5 bits of `rd`
	- 3 bits of `funct3`
	- 5 bits of `rs1`
	- 5 bits of `rs2`
	- 7 bitse of `funct7`
*/
func BuildInstructionR(opcode uint, rd uint, funct3 uint, rs1 uint, rs2 uint, funct7 uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, rd)
	builder.AddNextXBits(3, funct3)
	builder.AddNextXBits(5, rs1)
	builder.AddNextXBits(5, rs2)
	builder.AddNextXBits(7, funct7)
	return uint32(builder.Build())
}

/*BuildInstructionJ builds a 32-bit J instruction out of the arguments.
Uses lowest bits of arguments as follows:
	- 7 bits of `opcode`
	- 5 bits of `rd`
	- 20 bits of `offset`*/
func BuildInstructionJ(opcode uint, rd uint, offset uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, rd)
	builder.AddNextXBits(8, uint(Utils.GetBitsInInclusiveRange(offset, 11, 18)))
	builder.AddNextXBits(1, uint(Utils.GetBitsInInclusiveRange(offset, 10, 10)))
	builder.AddNextXBits(10, uint(Utils.GetBitsInInclusiveRange(offset, 0, 9)))
	builder.AddNextXBits(1, uint(Utils.GetBitsInInclusiveRange(offset, 19, 19)))

	return uint32(builder.Build())
}

/*BuildInstructionB builds a 32-bit B instruction out of the arguments.
Uses the lowest bits of arguments as follows:
	- 7 bits of `opcode`
	- 3 bits of `funct3`
	- 5 bits of `rs1`
	- 5 bits of `rs2`
	- 12 bits of `immediate`
*/
func BuildInstructionB(opcode uint, funct3 uint, rs1 uint, rs2 uint, immediate uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(1, uint(Utils.GetBitsInInclusiveRange(immediate, 10, 10)))
	builder.AddNextXBits(4, uint(Utils.GetBitsInInclusiveRange(immediate, 0, 3)))
	builder.AddNextXBits(3, funct3)
	builder.AddNextXBits(5, rs1)
	builder.AddNextXBits(5, rs2)
	builder.AddNextXBits(6, uint(Utils.GetBitsInInclusiveRange(immediate, 4, 9)))
	builder.AddNextXBits(1, uint(Utils.GetBitsInInclusiveRange(immediate, 11, 11)))

	return uint32(builder.Build())
}

func BuildInstructionS(opcode uint, funct3 uint, rs1 uint, rs2 uint, immediate uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, uint(Utils.GetBitsInInclusiveRange(immediate, 0, 4)))
	builder.AddNextXBits(3, funct3)
	builder.AddNextXBits(5, rs1)
	builder.AddNextXBits(5, rs2)
	builder.AddNextXBits(7, uint(Utils.GetBitsInInclusiveRange(immediate, 5, 11)))

	return uint32(builder.Build())
}
