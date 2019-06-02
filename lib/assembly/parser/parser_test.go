package parser

import (
	"testing"

	Assembler "github.com/chenhowa/computer/lib/assembly"
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

func (suite *RiscVParserSuite) TestAnd() {

}

type MockTokenStream struct {
	tokens       []MockToken
	currentToken uint
}

func makeMockTokenStream(tokens []MockToken) MockTokenStream {
	stream := MockTokenStream{
		tokens:       tokens,
		currentToken: 0,
	}
	return stream
}

type MockToken struct {
	tokenType         Assembler.TokenType
	tokenString       string
	charsSinceNewline Assembler.CharCount
}

func makeMockToken(tokenType Assembler.TokenType, tokenString string, count Assembler.CharCount) MockToken {
	token := MockToken{
		tokenType:         tokenType,
		tokenString:       tokenString,
		charsSinceNewline: count,
	}
	return token
}
