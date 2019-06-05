package parser

import (
	"errors"
	"fmt"
	"strings"

	Assembler "github.com/chenhowa/computer/lib/assembly"
)

/*RiscVParser is responsible for parsing Tokens that represent valid syntax in the RISC-V 32I assembly language, and
organizing them into an abstract tree for later code generation*/
type RiscVParser struct {
	lineCount Assembler.LineCount
}

/*MakeRiscVParser is a constructor for RiscVParser*/
func MakeRiscVParser() RiscVParser {
	parser := RiscVParser{
		lineCount: 0,
	}
	return parser
}

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
	GetAstNode() astNode
	GetParentIterator() (AstIterator, error)
	GetChildIterator(index uint) (AstIterator, error)
}

type astNode interface {
	GetLineCount() Assembler.LineCount
	GetCharCountSinceNewline() Assembler.CharCount
	GetTokenType() Assembler.TokenType
	GetTokenString() string
}

/*Parse takes a `tokenStream` and attempts to parse all the tokens in the stream into an Abstract Syntax Tree representation
of the RISC-V Assembly Program. If the parse is unsuccessful, it will return a non-nil error `err`.
If the parse is successful, it will return the AST `tree`, as well as `linesEncountered`, which represents the number
of newline tokens that were encountered in the parsing of `tokenStream` */
func (parser *RiscVParser) Parse(tokenStream tokenStream) (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {

	//optionalNewlines() && optionalInstructions() && optionalNewlines()

	// If no parses succeeded at all, all we can say is that the program could not be parsed
	ast := RiscVAst{
		root: RiscVAstNode{
			parent:    nil,
			lineCount: 0,
			//data:     ,
			children: []RiscVAstNode{},
		},
	}
	return ast, 0, errors.New("Parse: Input program could not be parsed at all")
}

/*
func optionalInstructions() (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {
	//noInstructions() || instructions()
}

func noInstructions() (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {

}

func instructions() (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {
	//instruction() && (noInstructions() || (newline() && instructions()))
}*/

/*RiscVAst represents an Abstract Syntax Tree of a valid RISC-V 32I Assembly Program*/
type RiscVAst struct {
	root RiscVAstNode
}

/*GetRootIterator returns an iterator to the root node of this AST*/
func (ast *RiscVAst) GetRootIterator() AstIterator {
	iter := makeRiscVAstIterator(&ast.root)
	return &iter
}

/*String returns the string representation of the AST*/
func (ast *RiscVAst) String() string {
	root := ast.GetRootIterator()
	str := convertToString(root)
	return str
}

func convertToString(iter AstIterator) string {
	var builder strings.Builder
	builder.WriteString("(" + iter.GetAstNode().GetTokenString() + ")")
	for i := uint(0); i < iter.GetNumChildren(); i++ {
		citer, err := iter.GetChildIterator(i)
		if err != nil {
			panic("convertToString: grabbed an invalid child iterator")
		}
		builder.WriteString(convertToString(citer))
	}

	return builder.String()
}

/*RiscVAstIterator is an object that aids iteration over an AST contianing Risc-V 32I Assembly Language tokens*/
type RiscVAstIterator struct {
	node *RiscVAstNode
}

func makeRiscVAstIterator(node *RiscVAstNode) RiscVAstIterator {
	iter := RiscVAstIterator{
		node: node,
	}
	return iter
}

/*GetNumChildren returns the number of children of the node pointed to by this iterator*/
func (it *RiscVAstIterator) GetNumChildren() uint {
	return uint(len(it.node.children))
}

/*GetAstNode returns a pointer to the naked node that this iterator points at*/
func (it *RiscVAstIterator) GetAstNode() astNode {
	return it.node
}

/*GetParentIterator returns an iterator to this iterator's node's parent, if it has one
Otherwise it returns an error*/
func (it *RiscVAstIterator) GetParentIterator() (AstIterator, error) {
	if it.node.parent != nil {
		iter := makeRiscVAstIterator(it.node.parent)
		return &iter, nil
	} else {
		iter := RiscVAstIterator{}
		return &iter, errors.New("GetParentIterator: No parent iterator")
	}
}

/*GetChildIterator returns an iterator to the i'th child of this iterator's node, if it has one.
Otherwise it returns an error*/
func (it *RiscVAstIterator) GetChildIterator(index uint) (AstIterator, error) {
	if index < uint(len(it.node.children)) {
		iter := makeRiscVAstIterator(&it.node.children[index])
		return &iter, nil
	} else {
		iter := RiscVAstIterator{}
		return &iter, fmt.Errorf("GetChildIterator: No child iterator with index %d", index)
	}
}

/*RiscVAstNode is a node within an AST of Risc-V 32I Assembly Language tokens
 */
type RiscVAstNode struct {
	parent    *RiscVAstNode
	lineCount Assembler.LineCount
	data      Token
	children  []RiscVAstNode
}

/*GetLineCount returns the program line of this particular node of the program*/
func (node *RiscVAstNode) GetLineCount() Assembler.LineCount {
	return node.lineCount
}

/*GetCharCountSinceNewline returns the number of chars since the newline for this
particular node's token*/
func (node *RiscVAstNode) GetCharCountSinceNewline() Assembler.CharCount {
	return node.data.GetCharCountSinceNewline()
}

/*GetTokenType returns the type of the token within this node*/
func (node *RiscVAstNode) GetTokenType() Assembler.TokenType {
	return node.data.GetTokenType()
}

/*GetTokenString returnst he string of the token within this node*/
func (node *RiscVAstNode) GetTokenString() string {
	return node.data.GetTokenString()
}
