package parser

import (
	"testing"

	Assembler "github.com/chenhowa/computer/lib/assembly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserUtilsSuite struct {
	suite.Suite
}

func TestParserUtilsSuite(t *testing.T) {
	suite.Run(t, new(ParserUtilsSuite))
}

func (suite *ParserUtilsSuite) SetupTest() {
}

func (suite *ParserUtilsSuite) TestMnemonic() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.AND, "AND", 0),
	}
	stream := makeMockTokenStream(tokens)

	ast, ok := mnemonic(&stream)
	assert.Equal(true, ok)
	assert.Equal(ast.String(), "(AND)")
}

func (suite *ParserUtilsSuite) TestOperand() {
	assert := assert.New(suite.T())
	var tokens = []MockToken{
		makeMockToken(Assembler.X1, "x1", 0),
	}
	stream := makeMockTokenStream(tokens)

	ast, ok := operand(&stream)
	assert.Equal(true, ok)
	assert.Equal(ast.String(), "(x1)")
}

func (suite *ParserUtilsSuite) TestIsOperand() {
	assert := assert.New(suite.T())
	token := makeMockToken(Assembler.X1, "x1", 0)
	assert.Equal(true, isOperand(&token))
}
