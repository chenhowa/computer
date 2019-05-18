package instructionParsing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParseSuite struct {
	suite.Suite
	parser  RiscVBinaryInstructionParser
	builder BinaryInstructionBuilder
}

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

func (suite *ParseSuite) SetupTest() {
	suite.parser = RiscVBinaryInstructionParser{}
	suite.builder = MakeInstructionBuilder(32)
}

func (suite *ParseSuite) TestParseImmediateArithmetic() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(ImmArith)) //opcode
	suite.builder.AddNextXBits(5, 3)              // rd (dest)
	suite.builder.AddNextXBits(3, 5)              // funct3
	suite.builder.AddNextXBits(5, 2)              // rs1 (src)
	suite.builder.AddNextXBits(12, 64)            // imm (immediate)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    I,
		OpCode:             ImmArith,
		FiveBitRegister1:   2,
		FiveBitDestination: 3,
		Funct3:             5,
		TwelveBitImmediate: 64,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseLUI() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(LUI)) // opcode
	suite.builder.AddNextXBits(5, 4)         // rd (dest)
	suite.builder.AddNextXBits(20, 128)      // imm (immediate)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    U,
		OpCode:             LUI,
		FiveBitDestination: 4,
		TwentyBitImmediate: 128,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseAUIPC() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(AUIPC)) // opcode
	suite.builder.AddNextXBits(5, 4)           // rd (dest)
	suite.builder.AddNextXBits(20, 128)        // imm (immediate)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    U,
		OpCode:             AUIPC,
		FiveBitDestination: 4,
		TwentyBitImmediate: 128,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseRegArith() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(RegArith)) // opcode
	suite.builder.AddNextXBits(5, 31)             // rd (dest)
	suite.builder.AddNextXBits(3, 7)              // funct3
	suite.builder.AddNextXBits(5, 1)              // rs1 (src1)
	suite.builder.AddNextXBits(5, 2)              // rs2 (src2)
	suite.builder.AddNextXBits(7, 18)             // funct7

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    R,
		OpCode:             RegArith,
		FiveBitDestination: 31,
		Funct3:             7,
		FiveBitRegister1:   1,
		FiveBitRegister2:   2,
		Funct7:             18,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseJAL() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(JAL)) // opcode
	suite.builder.AddNextXBits(5, 0)         // rd (dest)
	suite.builder.AddNextXBits(8, 0)         // imm[19:12]
	suite.builder.AddNextXBits(1, 1)         // imm[11]
	suite.builder.AddNextXBits(10, 0)        // imm[10:1]
	suite.builder.AddNextXBits(1, 1)         // imm[20]

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    J,
		OpCode:             JAL,
		FiveBitDestination: 0,
		TwentyBitImmediate: 525312,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseJALR() {
	assert := assert.New(suite.T())

	suite.builder.AddNextXBits(7, uint(JALR)) // opcode
	suite.builder.AddNextXBits(5, 3)          // rd (dest)
	suite.builder.AddNextXBits(3, 2)          // funct3
	suite.builder.AddNextXBits(5, 4)          // rs1 (base)
	suite.builder.AddNextXBits(12, 4095)      // imm (offset)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    I,
		OpCode:             JALR,
		FiveBitDestination: 3,
		Funct3:             2,
		FiveBitRegister1:   4,
		TwelveBitImmediate: 4095,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseBranch() {
	assert := assert.New(suite.T())
	builder := &suite.builder

	builder.AddNextXBits(7, uint(Branch)) // opcode
	builder.AddNextXBits(1, 1)            // imm[11]
	builder.AddNextXBits(4, 0)            // imm[4:1]
	builder.AddNextXBits(3, 7)            // funct3
	builder.AddNextXBits(5, 31)           // rs1 (src1)
	builder.AddNextXBits(5, 22)           // rs2 (src2)
	builder.AddNextXBits(6, 0)            // imm[10:5]
	builder.AddNextXBits(1, 1)            // imm[12]

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    B,
		OpCode:             Branch,
		Funct3:             7,
		FiveBitRegister1:   31,
		FiveBitRegister2:   22,
		TwelveBitImmediate: 3072,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseLoad() {
	assert := assert.New(suite.T())
	builder := &suite.builder

	builder.AddNextXBits(7, uint(Load))  // opcode
	suite.builder.AddNextXBits(5, 3)     // rd (dest)
	suite.builder.AddNextXBits(3, 2)     // funct3
	suite.builder.AddNextXBits(5, 4)     // rs1 (base)
	suite.builder.AddNextXBits(12, 4095) // imm (offset)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    I,
		OpCode:             Load,
		FiveBitDestination: 3,
		Funct3:             2,
		FiveBitRegister1:   4,
		TwelveBitImmediate: 4095,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseStore() {
	assert := assert.New(suite.T())
	builder := &suite.builder

	builder.AddNextXBits(7, uint(Store)) // opcode
	builder.AddNextXBits(5, 31)          // imm[4:0]
	builder.AddNextXBits(3, 1)           // funct3
	builder.AddNextXBits(5, 1)           // rs1 (base)
	builder.AddNextXBits(5, 4)           //rs2 (src)
	builder.AddNextXBits(7, 127)         //imm[11:5]

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    S,
		OpCode:             Store,
		Funct3:             1,
		FiveBitRegister1:   1,
		FiveBitRegister2:   4,
		TwelveBitImmediate: 4095,
	}

	assert.Equal(expected, actual)
}

func (suite *ParseSuite) TestParseSystem() {
	assert := assert.New(suite.T())
	builder := &suite.builder

	builder.AddNextXBits(7, uint(System)) // opcode
	suite.builder.AddNextXBits(5, 3)      // rd (dest)
	suite.builder.AddNextXBits(3, 2)      // funct3
	suite.builder.AddNextXBits(5, 4)      // rs1 (base)
	suite.builder.AddNextXBits(12, 3000)  // imm (offset)

	actual := suite.parser.Parse(uint32(suite.builder.Build()))

	expected := RiscVBinaryParseResult{
		InstructionType:    I,
		OpCode:             System,
		FiveBitDestination: 3,
		Funct3:             2,
		FiveBitRegister1:   4,
		TwelveBitImmediate: 3000,
	}

	assert.Equal(expected, actual)
}
