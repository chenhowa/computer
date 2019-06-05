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

func (suite *RiscVParserSuite) TestAnd() {

}
