package tokenizer

import (
	"fmt"
	"strings"

	Assembler "github.com/chenhowa/computer/lib/assembly"
)

/*RiscVTokenStream generates a stream of tokens from a string containing valid RISC-V assembly instructions*/
type RiscVTokenStream struct {
	input                 string
	currentPosition       uint
	charsSinceLastNewline uint
}

/*MakeRiscVTokenStream constructs a RiscVTokenStream*/
func MakeRiscVTokenStream(tokens string) RiscVTokenStream {
	stream := RiscVTokenStream{
		input:                 tokens,
		currentPosition:       0,
		charsSinceLastNewline: 0,
	}

	return stream
}

/*Next returns a the next token in the input stream*/
func (s *RiscVTokenStream) Next() (RiscVToken, error) {
	token, charsRead, err := s.getNextToken()
	s.updateCharsSinceLastNewline(&token, charsRead)
	if err == nil {
		s.charsSinceLastNewline += s.discardSkippableChars()
	}
	return token, err
}

func (s *RiscVTokenStream) hasMoreInput() bool {
	return s.currentPosition < uint(len(s.input))
}

func (s *RiscVTokenStream) discardSkippableChars() (charsDiscarded uint) {
	charCount := uint(0)

	discarded := s._discardSkippableChars()
	charCount += discarded
	for discarded > 0 {
		discarded = s._discardSkippableChars()
		charCount += discarded
	}
	return charCount
}

func (s *RiscVTokenStream) _discardSkippableChars() (charsDiscarded uint) {
	count := uint(0)

	/*discarding should be done in this order, to ensure that comments are removed */
	count += s.discardMultilineComment()
	count += s.discardSingleLineComment()
	count += s.discardSimpleChars()

	return count
}

func (s *RiscVTokenStream) discardSimpleChars() (charsDiscarded uint) {
	charCount := uint(0)
	char, err := s.getCurrentChar()
	for err == nil && !isUnskippableChar(char) {
		charCount++
		s.incrementCurrentPosition()
		char, err = s.getCurrentChar()
	}
	return charCount
}

func (s *RiscVTokenStream) discardSingleLineComment() (charsDiscarded uint) {
	charCount := uint(0)
	char, err := s.getCurrentChar()
	nextChar, err2 := s.getNextChar()

	if (err == nil && err2 == nil) && (char == '/' && nextChar == '/') {
		s.incrementCurrentPosition()
		s.incrementCurrentPosition()
		charCount += 2
		charCount += s.discardThroughSingleLineCommentEnd()
		return charCount
	} else {
		return 0
	}
}

func (s *RiscVTokenStream) discardThroughSingleLineCommentEnd() (charsDiscarded uint) {
	charCount := uint(0)
	char, err := s.getCurrentChar()
	for (err == nil) && (char != '\n') {
		charCount++
		s.incrementCurrentPosition()
		char, err = s.getCurrentChar()
	}

	return charCount
}

func (s *RiscVTokenStream) discardMultilineComment() (charsDiscarded uint) {
	charCount := uint(0)
	char, err := s.getCurrentChar()
	nextChar, err2 := s.getNextChar()

	if (err == nil && err2 == nil) && (char == '/' && nextChar == '*') {
		s.incrementCurrentPosition()
		s.incrementCurrentPosition()
		charCount += 2
		charCount += s.discardThroughMultilineCommentEnd()
		return charCount
	} else {
		return 0
	}
}

func (s *RiscVTokenStream) discardThroughMultilineCommentEnd() (charsDiscarded uint) {
	charCount := uint(0)
	char, err := s.getCurrentChar()
	nextChar, _ := s.getNextChar()

	for (err == nil) && (char != '*' || nextChar != '/') {
		charCount++
		s.incrementCurrentPosition()
		char, err = s.getCurrentChar()
		nextChar, _ = s.getNextChar()
	}

	if char == '*' && nextChar == '/' {
		s.incrementCurrentPosition()
		s.incrementCurrentPosition()
		charCount += 2
	} else {
		// we reached end of input. nothing doing.
	}

	return charCount
}

func (s *RiscVTokenStream) getNextToken() (RiscVToken, uint, error) {
	//Read one char of the next token at a time until the next skippable char is encountered.
	// That is the next token, which must be evaluated for being a valid token
	tokenString, charsRead := s.getNextTokenString()
	tokenType, err := getTokenType(tokenString)
	charsSinceLastNewline := Assembler.CharCount(s.charsSinceLastNewline + charsRead - uint(len(tokenString)))
	if err == nil {
		token := getCleanRiscVToken(tokenType, tokenString, charsSinceLastNewline)
		return token, charsRead, nil
	} else if tokenString == "" {
		return makeRiscVToken(Assembler.EndOfInput, tokenString, charsSinceLastNewline), charsRead, nil
	} else {
		return makeRiscVToken(Assembler.Error, tokenString, charsSinceLastNewline), charsRead, err
	}
}

func (s *RiscVTokenStream) getNextTokenString() (token string, charsRead uint) {
	charCount := uint(0)
	charCount += s.discardSkippableChars()
	_, err := s.getCurrentChar()
	if err != nil {
		return "", charCount
	} else {
		builder := strings.Builder{}
		char, err := s.getCurrentChar()
		for err == nil && continueReadingTokenInput(char, builder.String()) {
			charCount++
			builder.WriteByte(char)
			s.incrementCurrentPosition()
			char, err = s.getCurrentChar()
		}
		return builder.String(), charCount
	}
}

func getTokenType(tokenString string) (Assembler.TokenType, error) {
	tokenType, ok := mnemonicToToken[Mnemonic(tokenString)]
	if ok {
		return tokenType, nil
	}

	if tokenString == "\n" {
		return Assembler.Newline, nil
	}

	if isNumericConstant(tokenString) {
		return Assembler.NumericConstant, nil
	}

	if isLabel(tokenString) {
		return Assembler.Label, nil
	}

	return tokenType, fmt.Errorf("getTokenType: no token type found for this token %s", tokenString)
}

func (s *RiscVTokenStream) updateCharsSinceLastNewline(token *RiscVToken, charsRead uint) {
	if token.GetTokenType() == Assembler.Newline {
		s.charsSinceLastNewline = 0
	} else {
		s.charsSinceLastNewline += charsRead
	}
}

/*Save returns a tokenStreamReset. When the tokenStreamReset is invoked,
the tokenStream will be restored to its previous position in the input*/
func (s *RiscVTokenStream) Save() *RiscVTokenStreamReset {
	reset := makeRiscVTokenStreamReset(s.currentPosition, s.charsSinceLastNewline, s)
	return &reset
}

func (s *RiscVTokenStream) setCurrentInputPosition(position uint) {
	s.currentPosition = position
}

func (s *RiscVTokenStream) setCharsSinceLastNewline(count uint) {
	s.charsSinceLastNewline = count
}

/*RiscVToken is a token that holds the token type, the token's string,
and the token's character count since the last newline encountered
in the input string*/
type RiscVToken struct {
	tokenType             Assembler.TokenType
	token                 string
	charCountSinceNewline Assembler.CharCount
}

func makeRiscVToken(tokenType Assembler.TokenType, token string, count Assembler.CharCount) RiscVToken {
	riscVToken := RiscVToken{
		tokenType:             tokenType,
		token:                 token,
		charCountSinceNewline: count,
	}
	return riscVToken
}

/*GetTokenType returns the type of the token*/
func (token *RiscVToken) GetTokenType() Assembler.TokenType {
	return token.tokenType
}

/*GetTokenString returns the string that makes up this token*/
func (token *RiscVToken) GetTokenString() string {
	return token.token
}

/*GetCharCountSinceNewline returns the number of chars from the last newline
to the start of this token. This is to improve error messages to help locate
the offending token*/
func (token *RiscVToken) GetCharCountSinceNewline() Assembler.CharCount {
	return token.charCountSinceNewline
}

/*RiscVTokenStreamReset is a Command Pattern object that contains the ability to reset
the token stream that produced it to a given position in the input*/
type RiscVTokenStreamReset struct {
	inputPosition         uint
	charsSinceLastNewline uint
	stream                *RiscVTokenStream
}

func makeRiscVTokenStreamReset(inputPosition uint, charsSinceLastNewline uint, stream *RiscVTokenStream) RiscVTokenStreamReset {
	reset := RiscVTokenStreamReset{
		inputPosition:         inputPosition,
		stream:                stream,
		charsSinceLastNewline: charsSinceLastNewline,
	}
	return reset
}

/*Reset is an sets the token stream that produced this RiscVTokenStreamReset to be
at the input position contained within this RiscVTokenStreamReset*/
func (r *RiscVTokenStreamReset) Reset() {
	r.stream.setCurrentInputPosition(r.inputPosition)
	r.stream.setCharsSinceLastNewline(r.charsSinceLastNewline)
}
