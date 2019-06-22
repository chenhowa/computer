package parser

import (
	"fmt"
	"testing"

	Assembler "github.com/chenhowa/computer/lib/assembly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RiscVParserSuite struct {
	suite.Suite
	parser *RiscVParser
}

func TestRiscVParserSuite(t *testing.T) {
	suite.Run(t, new(RiscVParserSuite))
}

func (suite *RiscVParserSuite) SetupTest() {
	parser := MakeRiscVParser()
	suite.parser = &parser
}

func (suite *RiscVParserSuite) TestSanity_OneFakeToken() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.Error, "ERROR", 0),
	}
	stream := makeMockTokenStream(tokens)
	_, _, err := suite.parser.Parse(&stream)
	assert.Equal("Parse: Input program could not be parsed at all", err.Error())
}

func (suite *RiscVParserSuite) TestNewlinesParsing() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.Newline, "NEW", 0),
		makeMockToken(Assembler.Newline, "NEW", 0),
		makeMockToken(Assembler.Newline, "NEW", 0),
	}
	stream := makeMockTokenStream(tokens)
	ast, linecount, err := suite.parser.Parse(&stream)
	assert.Equal(nil, err)
	assert.Equal(Assembler.LineCount(3), linecount)
	assert.Equal("(NEW(NEW(NEW)))", ast.String())
}

func (suite *RiscVParserSuite) TestAnd() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.X1, "x1", 1),
		makeMockToken(Assembler.X2, "x2", 2),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(AND(x1)(x2))))", Assembler.Instructions, Assembler.Instruction), ast.String())
}

func (suite *RiscVParserSuite) TestAnd_Multiple() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.X1, "x1", 1),
		makeMockToken(Assembler.X2, "x2", 2),
		makeMockToken(Assembler.Newline, "NEW", 0),
		makeMockToken(Assembler.Newline, "NEW", 0),
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.X3, "x3", 1),
		makeMockToken(Assembler.X4, "x4", 2),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(2), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(AND(x1)(x2)))(1(AND(x3)(x4))))", Assembler.Instructions, Assembler.Instruction), ast.String())
}

func (suite *RiscVParserSuite) TestLabel() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.Label, "START_LOOP", 0),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(START_LOOP)))", Assembler.Instructions, Assembler.Instruction), ast.String())
}

func (suite *RiscVParserSuite) TestIdentifier() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.Identifier, "Yolo", 1),
		makeMockToken(Assembler.Identifier, "Cholo", 2),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(AND(Yolo)(Cholo))))", Assembler.Instructions, Assembler.Instruction), ast.String())

}

func (suite *RiscVParserSuite) TestNumericConstants() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.NumericConstant, "1000", 1),
		makeMockToken(Assembler.NumericConstant, "30,000", 2),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(AND(1000)(30,000))))", Assembler.Instructions, Assembler.Instruction), ast.String())
}

func (suite *RiscVParserSuite) TestRegisterAndImmediate() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.RegisterAndImmediate, "10(x1)", 1),
		makeMockToken(Assembler.RegisterAndImmediate, "20(x2)", 2),
	}
	stream := makeMockTokenStream(tokens)
	ast, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	assert.Equal(fmt.Sprintf("(%d(%d(AND(10(x1))(20(x2)))))", Assembler.Instructions, Assembler.Instruction), ast.String())
}
