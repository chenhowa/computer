package parser

import (
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

/*
func (suite *RiscVParserSuite) TestAnd() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
		makeMockToken(Assembler.X1, "x1", 1),
		makeMockToken(Assembler.X2, "x2", 2),
	}
	stream := makeMockTokenStream(tokens)
	_, lines, err := suite.parser.Parse(&stream)
	assert.Equal(uint(0), uint(lines))
	assert.Equal(nil, err)
	//assert.Equal(ast.String(), "(AND(x1)(x2))")
}
*/
