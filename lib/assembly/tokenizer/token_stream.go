package tokenizer

import Assembler "github.com/chenhowa/computer/lib/assembly"

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
func (s *RiscVTokenStream) Next() RiscVToken {

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
	tokenType Assembler.TokenType
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
