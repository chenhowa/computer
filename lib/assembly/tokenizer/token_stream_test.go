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

func (suite *RiscVTokenStreamSuite) TestNext_Labels_Success() {
	//Need specific format for validity: Capitalized word followed immediately by a colon.
	input := "Else: ADDI"
	stream := MakeRiscVTokenStream(input)

	expected := makeRiscVToken(Assembler.Label, string("Else"), Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), Assembler.CharCount(uint(len("Else: "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Labels_Failure() {
	input := "Else:ADDI"
	stream := MakeRiscVTokenStream(input)
	_, err := stream.Next()
	assert.NotEqual(suite.T(), nil, err)

	input = "ELse:"
	stream = MakeRiscVTokenStream(input)
	_, err = stream.Next()
	assert.NotEqual(suite.T(), nil, err)

	input = "else:"
	stream = MakeRiscVTokenStream(input)
	_, err = stream.Next()
	assert.NotEqual(suite.T(), nil, err)
}

func (suite *RiscVTokenStreamSuite) TestNext_Registers_Success() {
	input := "x0 ADDI"
	stream := MakeRiscVTokenStream(input)

	expected := makeRiscVToken(Assembler.X0, "x0", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), Assembler.CharCount(uint(len("x0 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "x31 ADDI"
	stream = MakeRiscVTokenStream(input)

	expected = makeRiscVToken(Assembler.X31, "x31", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), Assembler.CharCount(uint(len("x31 "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Registers_Failure() {
	input := "x32"
	stream := MakeRiscVTokenStream(input)
	_, err := stream.Next()
	assert.NotEqual(suite.T(), nil, err)
}

func (suite *RiscVTokenStreamSuite) TestNext_Memory_RegisterImmediatePair_Success() {
	input := "8(x0)"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.RegisterAndImmediate, "8(x0)", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	input = "1000(x31) ADDI"
	stream = MakeRiscVTokenStream(input)
	expected = makeRiscVToken(Assembler.RegisterAndImmediate, "1000(x31)", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.ADDI, string(ADDI), Assembler.CharCount(uint(len("1000(x31) "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Memory_RegisterImmediatePair_Failure() {
	input := "8(x32)"
	stream := MakeRiscVTokenStream(input)
	_, err := stream.Next()
	assert.NotEqual(suite.T(), nil, err)

	input = "8(x31)ADDI"
	stream = MakeRiscVTokenStream(input)
	_, err = stream.Next()
	assert.NotEqual(suite.T(), nil, err)
}

func (suite *RiscVTokenStreamSuite) TestNext_RegisterNickname_1() {
	input := "zero ra sp gp tp t0 t1 t2 s0 fp s1"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.X0, "zero", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X1, "ra", Assembler.CharCount(uint(len("zero "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X2, "sp", Assembler.CharCount(uint(len("zero ra "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X3, "gp", Assembler.CharCount(uint(len("zero ra sp "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X4, "tp", Assembler.CharCount(uint(len("zero ra sp gp "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X5, "t0", Assembler.CharCount(uint(len("zero ra sp gp tp "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X6, "t1", Assembler.CharCount(uint(len("zero ra sp gp tp t0 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X7, "t2", Assembler.CharCount(uint(len("zero ra sp gp tp t0 t1 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X8, "s0", Assembler.CharCount(uint(len("zero ra sp gp tp t0 t1 t2 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X8, "fp", Assembler.CharCount(uint(len("zero ra sp gp tp t0 t1 t2 s0 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X9, "s1", Assembler.CharCount(uint(len("zero ra sp gp tp t0 t1 t2 s0 fp "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_RegisterNickname_2() {
	input := "a0 a1 a2 a3 a4 a5 a6 a7"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.X10, "a0", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X11, "a1", Assembler.CharCount(uint(len("a0 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X12, "a2", Assembler.CharCount(uint(len("a0 a1 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X13, "a3", Assembler.CharCount(uint(len("a0 a1 a2 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X14, "a4", Assembler.CharCount(uint(len("a0 a1 a2 a3 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X15, "a5", Assembler.CharCount(uint(len("a0 a1 a2 a3 a4 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X16, "a6", Assembler.CharCount(uint(len("a0 a1 a2 a3 a4 a5 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X17, "a7", Assembler.CharCount(uint(len("a0 a1 a2 a3 a4 a5 a6 "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_RegisterNickname_3() {
	input := "s2 s3 s4 s5 s6 s7 s8 s9 s10 s11"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.X18, "s2", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X19, "s3", Assembler.CharCount(uint(len("s2 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X20, "s4", Assembler.CharCount(uint(len("s2 s3 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X21, "s5", Assembler.CharCount(uint(len("s2 s3 s4 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X22, "s6", Assembler.CharCount(uint(len("s2 s3 s4 s5 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X23, "s7", Assembler.CharCount(uint(len("s2 s3 s4 s5 s6 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X24, "s8", Assembler.CharCount(uint(len("s2 s3 s4 s5 s6 s7 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X25, "s9", Assembler.CharCount(uint(len("s2 s3 s4 s5 s6 s7 s8 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X26, "s10", Assembler.CharCount(uint(len("s2 s3 s4 s5 s6 s7 s8 s9 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X27, "s11", Assembler.CharCount(uint(len("s2 s3 s4 s5 s6 s7 s8 s9 s10 "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_RegisterNickname_4() {
	input := "t3 t4 t5 t6"
	stream := MakeRiscVTokenStream(input)
	expected := makeRiscVToken(Assembler.X28, "t3", Assembler.CharCount(0))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X29, "t4", Assembler.CharCount(uint(len("t3 "))))
	suite.AssertNextTokenIs(&stream, &expected)

	expected = makeRiscVToken(Assembler.X30, "t5", Assembler.CharCount(uint(len("t3 t4 "))))
	suite.AssertNextTokenIs(&stream, &expected)
	expected = makeRiscVToken(Assembler.X31, "t6", Assembler.CharCount(uint(len("t3 t4 t5 "))))
	suite.AssertNextTokenIs(&stream, &expected)
}

func (suite *RiscVTokenStreamSuite) TestNext_Identifier() {

}

func (suite *RiscVTokenStreamSuite) Test_ReadsThroughMultipleErrors() {
	/*This test proves the token stream will continue to consume input even when
	having read an unrecognized token. This way the caller can report multiple errors
	if necessary. Or should there be a TokenType for unrecognized tokens? Since Error implies
	that something serious has occurred ... - but error allows a string. So the semantics
	of this function shoudl be that 'error' is recoverable, which is exactly what I want.
	So no need for new tokentype.*/
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
