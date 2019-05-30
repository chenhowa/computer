package parser

import (
	Assembler "github.com/chenhowa/computer/lib/assembly"
)

/*RiscVParser is responsible for parsing Tokens that represent valid syntax in the RISC-V 32I assembly language, and
organizing them into an abstract tree for later code generation*/
type RiscVParser struct {
}

type tokenStream interface {
	HasNext() bool
	Next() (token, error)
	Save() tokenStreamReset
}

type tokenStreamReset interface {
	Reset()
}

type token interface {
	GetTokenType() Assembler.TokenType
	GetTokenString() string
	GetCharCountSinceNewline() Assembler.CharCount
}

type abstractSyntaxTree interface {
	GetRootIterator() astIterator
}

type astIterator interface {
	GetNumChildren() uint
	GetAstNode() astNode
	GetParentIterator() (astIterator, error)
	GetChildIterator(index uint) (astIterator, error)
}

type astNode interface {
	GetLineCount()
	GetCharCountSinceNewline() Assembler.CharCount
	GetTokenType() Assembler.TokenType
	GetTokenString() string
}

/*Parse takes a `tokenStream` and a `lineCount` and attempts to parse all the tokens in the stream into an Abstract Syntax Tree representation
of the RISC-V Assembly Program. If the parse is unsuccessful, it will return a non-nil error `err`
If the parse is successful, it will return the AST `tree`, as well as a `newLineCount` that represents the number
of newline tokens that was encountered in the `tokenStream` */
func (parser *RiscVParser) Parse(tokenStream tokenStream) (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {

}

/*RiscVAst represents an Abstract Syntax Tree of a valid RISC-V 32I Assembly Program*/
type RiscVAst struct{}
