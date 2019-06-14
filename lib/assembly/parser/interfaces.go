package parser

import Assembler "github.com/chenhowa/computer/lib/assembly"

type tokenStream interface {
	HasNext() bool
	Next() (Token, error)
	Save() TokenStreamReset
}

/*TokenStreamReset is an interface that can reset a tokenStream
 */
type TokenStreamReset interface {
	Reset()
}

/*Token is an interface that represents a token of the Risc-V 32I Assembly Language*/
type Token interface {
	GetTokenType() Assembler.TokenType
	GetTokenString() string
	GetCharCountSinceNewline() Assembler.CharCount
}

type abstractSyntaxTree interface {
	GetRootIterator() AstIterator
}

/*AstIterator is an interface that represents an iterator over an AST of RiscV Tokens*/
type AstIterator interface {
	GetNumChildren() uint
	GetAstNode() AstNode
	GetParentIterator() (AstIterator, error)
	GetChildIterator(index uint) (AstIterator, error)
}

/*AstNode is an interface that represents a node in an AST of RiscV Tokens*/
type AstNode interface {
	GetLineCount() Assembler.LineCount
	GetCharCountSinceNewline() Assembler.CharCount
	GetTokenType() Assembler.TokenType
	GetTokenString() string
}
