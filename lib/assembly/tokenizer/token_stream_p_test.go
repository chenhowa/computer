package tokenizer

import (
	Assembler "github.com/chenhowa/computer/lib/assembly"
	"github.com/stretchr/testify/assert"
)

func (suite *RiscVTokenStreamSuite) TestPrivate_getCurrentChar() {
	assert := assert.New(suite.T())
	stream := MakeRiscVTokenStream(string(ADDI))
	b, e := stream.getCurrentChar()
	assert.Equal(nil, e)
	assert.Equal(byte('A'), b)
}

func (suite *RiscVTokenStreamSuite) TestPrivate_discardSkippableChars() {
	assert := assert.New(suite.T())
	stream := MakeRiscVTokenStream(string(ADDI))
	discardCount := stream.discardSkippableChars()
	assert.Equal(uint(0), discardCount)
}

func (suite *RiscVTokenStreamSuite) TestPrivate_getTokenType() {
	assert := assert.New(suite.T())
	tokenType, err := getTokenType(string(ADDI))
	assert.Equal(nil, err)
	assert.Equal(Assembler.ADDI, tokenType)
}

func (suite *RiscVTokenStreamSuite) TestPrivate_getNextTokenString() {
	assert := assert.New(suite.T())
	stream := MakeRiscVTokenStream(string(ADDI))
	tokenString, charsRead := stream.getNextTokenString()
	assert.Equal(uint(len(string(ADDI))), charsRead)
	assert.Equal(string(ADDI), tokenString)
}
