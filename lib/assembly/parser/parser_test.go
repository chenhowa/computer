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

func (s *MockTokenStream) HasNext() bool {
	return s.currentToken < uint(len(s.tokens))
}

func (s *MockTokenStream) Next() (Token, error) {
	if s.HasNext() {
		token := s.tokens[s.currentToken]
		s.currentToken++
		return &token, nil
	} else {
		end := makeMockToken(Assembler.EndOfInput, "", 0)
		return &end, nil
	}
}

func (s *MockTokenStream) setPosition(position uint) {
	s.currentToken = position
}

func (s *MockTokenStream) Save() TokenStreamReset {
	reset := MockReset{
		position: s.currentToken,
		stream:   s,
	}

	return &reset
}

type MockReset struct {
	position uint
	stream   *MockTokenStream
}

func (r *MockReset) Reset() {
	r.stream.setPosition(r.position)
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

func (t *MockToken) GetTokenType() Assembler.TokenType {
	return t.tokenType
}

func (t *MockToken) GetTokenString() string {
	return t.tokenString
}

func (t *MockToken) GetCharCountSinceNewline() Assembler.CharCount {
	return t.charsSinceNewline
}

func (suite *RiscVParserSuite) TestSanity_OneToken() {
	var tokens = []MockToken{
		makeMockToken(Assembler.ADD, "ADD", 0),
	}
	stream := makeMockTokenStream(tokens)
	suite.parser.Parse(&stream)
}

func (suite *RiscVParserSuite) TestAnd() {

}
