package tokenizer

import (
	Assembler "github.com/chenhowa/computer/lib/assembly"
)

/*RiscVTokenizer is responsible for producing a token stream from the input string*/
type RiscVTokenizer struct {
}

type tokenStream interface {
	Next() token
	Save() tokenStreamReset
}

type token interface {
	GetTokenType() Assembler.TokenType
	GetTokenString() string
	GetCharCountSinceNewline() Assembler.CharCount
}

type tokenStreamReset interface {
	Reset()
}

/*Tokenize produces an object with interface `tokenStream`*/
func (t *RiscVTokenizer) Tokenize(tokens string) tokenStream {
	stream := MakeRiscVTokenStream(tokens)
	return stream
}
