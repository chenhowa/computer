package tokenizer

import (
	Assembler "github.com/chenhowa/computer/lib/assembly"
)

/*RiscVTokenStream generates a stream of tokens from a string containing valid RISC-V assembly instructions*/
type RiscVTokenStream struct {
	input           string
	currentPosition uint
}

/*MakeRiscVTokenStream constructs a RiscVTokenStream*/
func MakeRiscVTokenStream(tokens string) RiscVTokenStream {
	stream := RiscVTokenStream{
		input:           tokens,
		currentPosition: 0,
	}

	return stream
}

/*Next returns a the next token in the input stream*/
func (s *RiscVTokenStream) Next() (RiscVToken, error) {
	if s.hasMoreInput() {
		s.discardSkippableChars()
		token, err := s.getNextToken()
		if err == nil {
			s.discardSkippableChars()
		}
		return token, err
	} else {
		token := RiscVToken{
			tokenType: Assembler.EndOfInput,
		}
		return token, nil
	}
}

func (s *RiscVTokenStream) hasMoreInput() bool {
	return s.currentPosition < uint(len(s.input))
}

func (s *RiscVTokenStream) discardSkippableChars() {
	for char := s.input[s.currentPosition]; !isUnskippableChar(char); s.incrementCurrentPosition() {
	}
}

func isUnskippableChar(val byte) bool {
	return true
}

func (s *RiscVTokenStream) incrementCurrentPosition() {
	s.currentPosition++
}

func (s *RiscVTokenStream) getNextToken() (RiscVToken, error) {
	//Read one char of the next token at a time until the next skippable char is encountered.
	// That is the next token, which must be evaluated for being a valid token
}

/*Save returns a tokenStreamReset. When the tokenStreamReset is invoked,
the tokenStream will be restored to its previous position in the input*/
func (s *RiscVTokenStream) Save() *RiscVTokenStreamReset {
	reset := makeRiscVTokenStreamReset(s.currentPosition, s)
	return &reset
}

func (s *RiscVTokenStream) setCurrentInputPosition(position uint) {
	s.currentPosition = position
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
	inputPosition uint
	stream        *RiscVTokenStream
}

func makeRiscVTokenStreamReset(inputPosition uint, stream *RiscVTokenStream) RiscVTokenStreamReset {
	reset := RiscVTokenStreamReset{
		inputPosition: inputPosition,
		stream:        stream,
	}
	return reset
}

/*Reset is an sets the token stream that produced this RiscVTokenStreamReset to be
at the input position contained within this RiscVTokenStreamReset*/
func (r *RiscVTokenStreamReset) Reset() {
	r.stream.setCurrentInputPosition(r.inputPosition)
}
