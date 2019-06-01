package tokenizer

import (
	"fmt"
	"testing"

	Assembler "github.com/chenhowa/computer/lib/assembly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RiscVTokenStreamSuite struct {
	suite.Suite
	stream *RiscVTokenStream
}

func TestRiscVTokenStreamSuite(t *testing.T) {
	suite.Run(t, new(RiscVTokenStreamSuite))
}

func (suite *RiscVTokenStreamSuite) SetupTest() {
}

func (suite *RiscVTokenStreamSuite) AssertExpectedTokenEqualsActual(expected *RiscVToken, actual *RiscVToken) {
	assert := assert.New(suite.T())
	assert.Equal(expected.GetTokenType(), actual.GetTokenType(), "Token types should be the same")
	assert.Equal(expected.GetTokenString(), actual.GetTokenString(), "Token Strings should be the same")
	assert.Equal(expected.GetCharCountSinceNewline(), actual.GetCharCountSinceNewline(), "Characters since Newline should be the same")
}

func (suite *RiscVTokenStreamSuite) AssertNextTokenIs(stream *RiscVTokenStream, expected *RiscVToken) {
	assert := assert.New(suite.T())
	token, err := stream.Next()
	assert.Equal(nil, err)
	suite.AssertExpectedTokenEqualsActual(expected, &token)
}

func (suite *RiscVTokenStreamSuite) TestNext_Empty() {
	assert := assert.New(suite.T())
	stream := MakeRiscVTokenStream("")
	token, _ := stream.Next()
	expected := makeRiscVToken(Assembler.EndOfInput, "", 0)
	assert.Equal(expected.GetTokenType(), Assembler.EndOfInput)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)
}

func (suite *RiscVTokenStreamSuite) TestNext_OneMnemonicToken() {
	assert := assert.New(suite.T())
	stream := MakeRiscVTokenStream(string(ADDI))
	token, err := stream.Next()
	expected := makeRiscVToken(Assembler.ADDI, string(ADDI), 0)
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.ADDI)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)
}

func (suite *RiscVTokenStreamSuite) TestNext_TwoMnemonicTokens() {
	assert := assert.New(suite.T())
	input := fmt.Sprintf("%s %s", string(ADDI), string(SUB))
	stream := MakeRiscVTokenStream(input)
	token, err := stream.Next()
	expected := makeRiscVToken(Assembler.ADDI, string(ADDI), 0)
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.ADDI)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)

	token, err = stream.Next()
	expected = makeRiscVToken(Assembler.SUB, string(SUB), Assembler.CharCount(uint(len(ADDI)+1)))
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.SUB)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)
}

func (suite *RiscVTokenStreamSuite) TestNext_TwoMnemonicTokensAndNewline() {
	assert := assert.New(suite.T())
	input := fmt.Sprintf("%s \n %s", string(ADDI), string(SUB))
	stream := MakeRiscVTokenStream(input)

	token, err := stream.Next()
	expected := makeRiscVToken(Assembler.ADDI, string(ADDI), 0)
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.ADDI)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)

	token, err = stream.Next()
	expected = makeRiscVToken(Assembler.Newline, "\n", Assembler.CharCount(uint(len(ADDI)+1)))
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.Newline)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)

	token, err = stream.Next()
	expected = makeRiscVToken(Assembler.SUB, string(SUB), Assembler.CharCount(uint(1)))
	assert.Equal(nil, err)
	assert.Equal(expected.GetTokenType(), Assembler.SUB)
	suite.AssertExpectedTokenEqualsActual(&expected, &token)
}

func (suite *RiscVTokenStreamSuite) TestNext_NumericConstants() {
	input := "1 213 1,123 1,000,000"
	stream := MakeRiscVTokenStream(input)

	expected := makeRiscVToken(Assembler.NumericConstant, string("1"), Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.NumericConstant, string("213"), Assembler.CharCount(2))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.NumericConstant, string("1123"), Assembler.CharCount(6))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.NumericConstant, string("1000000"), Assembler.CharCount(12))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Labels() {
	//Need specific format for validity
	asdfdasf LABELS NEED TO LOOK UP RISC-V ASSEMBLY FOR THIS
}

func (suite *RiscVTokenStreamSuite) TestNext_Comments_Multiline() {
	stream := MakeRiscVTokenStream("/*")
	expected := makeRiscVToken(Assembler.EndOfInput, "", 2)
	suite.AssertNextTokenIs(&stream, &expected)

	stream = MakeRiscVTokenStream("/* asdfdas")
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(3+len("asdfdas"))))
	suite.AssertNextTokenIs(&stream, &expected)

	stream = MakeRiscVTokenStream("/* asdfdas\n*/")
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(3+len("asdfdas")+3)))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Comments_Multiline_With_Mnemonics() {
	input := "ADDI/**/"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.ADDI, string(ADDI), 0)
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len(input))))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "/**/ADDI"
	stream = MakeRiscVTokenStream(input)
	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), 4)
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len(input))))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "ADDI/**/ADDI"
	stream = MakeRiscVTokenStream(input)
	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), 0)
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), 8)
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len(input))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Comments_Singleline() {
	input := "//"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len(input))))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "// asdfdas"
	stream = MakeRiscVTokenStream(input)
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len(input))))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "// asdfdas\n\t \t"
	stream = MakeRiscVTokenStream(input)
	expected = makeRiscVToken(Assembler.Newline, "\n", Assembler.CharCount(uint(len("// asdfdas"))))
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.EndOfInput, "", Assembler.CharCount(uint(len("\t \t"))))
	suite.AssertNextTokenIs(&stream, &expected)
}
