package binaryInstructionExecution

import (
	Parser "github.com/chenhowa/computer/lib/binaryInstructionExecution/instructionParsing"
)

/*BuildInstructionI builds a 32-bit I Format instruction out of the arguments
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

/*BuildInstructionU builds a 32-bit U Format instruction out of the arguments*/
func BuildInstructionU(opcode uint, rd uint, immediate uint) uint32 {
	builder := Parser.MakeInstructionBuilder(32)
	builder.AddNextXBits(7, opcode)
	builder.AddNextXBits(5, rd)
	builder.AddNextXBits(20, immediate)
	return uint32(builder.Build())
}
