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

/*Parse takes a `tokenStream` and attempts to parse all the tokens in the stream into an Abstract Syntax Tree representation
of the RISC-V Assembly Program. If the parse is unsuccessful, it will return a non-nil error `err`.
If the parse is successful, it will return the AST `tree`, as well as `linesEncountered`, which represents the number
of newline tokens that were encountered in the parsing of `tokenStream` */
func (parser *RiscVParser) Parse(tokenStream tokenStream) (tree RiscVAst, linesEncountered Assembler.LineCount, err error) {

	//optionalNewlines() && optionalInstructions() && optionalNewlines()
	newlinesAst, newlinesOk := optionalNewlines(tokenStream)
	if newlinesOk {
		return newlinesAst, 0, nil
	}

	//Since it's optional, it doesn't matter whether it succeeded, or failed
	instructionsAst, instructionsOk := optionalInstructions(tokenStream)

	if instructionsOk {
		return instructionsAst, 0, nil
	}
	/*
		// Since it's optional, it doesn't matter whether it succeeded, or failed.
		newlinesAst2, newlinesOk2 := optionalNewlines(tokenStream)

		if newlinesOk || instructionsOk || newlinesOk2 {
			finalTree := combine(newlinesAst, instructionsAst, newlinesAst2)

			return finalTree, getLineCount(finalTree), nil
		}*/

	// If no parses succeeded at all, all we can say is that the program could not be parsed
	node := RiscVAstNode{
		parent:    nil,
		lineCount: 0,
		//data:     ,
		children: []*RiscVAstNode{},
	}
	errorAst := RiscVAst{
		root: &node,
	}
	return errorAst, 0, errors.New("Parse: Input program could not be parsed at all")
}

func optionalInstructions(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	optionalInstructionsAst, ok := _optionalInstructions(stream)

	if ok {
		return optionalInstructionsAst, true
	} else {
		reset.Reset()
		return RiscVAst{}, false
	}
}

func _optionalInstructions(stream tokenStream) (tree RiscVAst, success bool) {
	instructionAst, ok := instruction(stream)
	if ok {
		newlineAst, ok := newline(stream) // we require a newline between each instruction
		if ok {
			_optionalNewlines(stream) // more than 1 newline is acceptable, but not required.
			remAst, ok := _optionalInstructions(stream)
			if ok {

			} else {
				return instructionAst, true
			}
		} else {
			return instructionAst, true
		}

	} else {
		return RiscVAst{}, false
	}
}

func instruction(stream tokenStream) (tree RiscVAst, success bool) {
	mnemonicInstructionAst, mnemonicOk := mnemonicInstruction(stream)
	if mnemonicOk {
		return mnemonicInstructionAst, true
	}

	labelInstructionAst, labelOk := labelInstruction(stream)

	if labelOk {
		return labelInstructionAst, true
	}
}

func mnemonicInstruction(stream tokenStream) (tree RiscVAst, success bool) {
	mnemonicAst, mnemonicOk := mnemonic(stream)

	if !mnemonicOk {
		return RiscVAst{}, false
	}

	operand1Ast, operand1Ok := operand(stream)

	if !operand1Ok {
		return RiscVAst{}, false
	}

	operand2Ast, operand2Ok := operand(stream)

	if !operand2Ok {
		return RiscVAst{}, false
	}
}

func mnemonic(stream tokenStream) (tree RiscVAst, success bool) {
	token, tokenErr := stream.Next()

	if tokenErr == nil && isMnemonic(token.GetTokenType()) {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		ast := RiscVAst{}
		return ast, false
	}
}

func operand(stream tokenStream) (tree RiscVAst, success bool) {
	token, tokenErr := stream.Next()

	if tokenErr == nil && isOperand(token.GetTokenType()) {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		ast := RiscVAst{}
		return ast, false
	}
}

func labelInstruction(stream tokenStream) (tree RiscVAst, success bool) {
	labelAst, labelOk := label(stream)

	if labelOk {
		return labelAst, true
	} else {
		return RiscVAst{}, false
	}
}

func label(stream tokenStream) (tree RiscVAst, success bool) {
	token, tokenErr := stream.Next()

	if tokenErr == nil && isLabel(token.GetTokenType()) {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		ast := RiscVAst{}
		return ast, false
	}
}

func optionalNewlines(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	optionalNewlinesAst, ok := _optionalNewlines(stream)

	// Since this is optional newlines, as long as there was no stream error,
	// it is considered a successful parse.
	if ok {
		return optionalNewlinesAst, true
	} else {
		reset.Reset()
		return RiscVAst{}, false
	}
}

func _optionalNewlines(stream tokenStream) (tree RiscVAst, success bool) {
	newlineAst, ok := newline(stream)
	if ok {
		remAst, ok := _optionalNewlines(stream)
		if ok {
			parent := newlineAst.getRootIterator()
			return addAsChild(newlineAst, parent, remAst), true
		} else {
			return newlineAst, true
		}
	} else {
		return RiscVAst{}, false
	}
}

func addAsChild(parentTree RiscVAst, nodeToAddAt RiscVAstIterator, childTree RiscVAst) RiscVAst {
	tree := RiscVAst{
		root: parentTree.root,
	}

	nodeToAddAt.addAsNextChild(childTree.root)

	return tree
}

func newline(stream tokenStream) (tree RiscVAst, success bool) {
	token, tokenErr := stream.Next()
	if tokenErr == nil && token.GetTokenType() == Assembler.Newline {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		ast := RiscVAst{}
		return ast, false
	}
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
	root *RiscVAstNode
}

/*GetRootIterator returns an iterator to the root node of this AST*/
func (ast *RiscVAst) GetRootIterator() AstIterator {
	iter := ast.getRootIterator()
	return &iter
}

func (ast *RiscVAst) getRootIterator() RiscVAstIterator {
	iter := makeRiscVAstIterator(ast.root)
	return iter
}

/*String returns the string representation of the AST*/
func (ast *RiscVAst) String() string {
	root := ast.GetRootIterator()
	str := convertToString(root)
	return str
}

func convertToString(iter AstIterator) string {
	var builder strings.Builder
	builder.WriteString("(" + iter.GetAstNode().GetTokenString())
	for i := uint(0); i < iter.GetNumChildren(); i++ {
		citer, err := iter.GetChildIterator(i)
		if err != nil {
			panic("convertToString: grabbed an invalid child iterator")
		}
		builder.WriteString(convertToString(citer))
	}
	builder.WriteString(")")

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

func (it *RiscVAstIterator) addAsNextChild(node *RiscVAstNode) {
	it.node.children = append(it.node.children, node)
	node.parent = it.getAstNode()
}

/*GetNumChildren returns the number of children of the node pointed to by this iterator*/
func (it *RiscVAstIterator) GetNumChildren() uint {
	return uint(len(it.node.children))
}

/*GetAstNode returns a pointer to the naked node that this iterator points at*/
func (it *RiscVAstIterator) GetAstNode() AstNode {
	return it.getAstNode()
}

func (it *RiscVAstIterator) getAstNode() *RiscVAstNode {
	return it.node
}

/*GetParentIterator returns an iterator to this iterator's node's parent, if it has one
Otherwise it returns an error*/
func (it *RiscVAstIterator) GetParentIterator() (AstIterator, error) {
	return it.getParentIterator()
}

func (it *RiscVAstIterator) getParentIterator() (*RiscVAstIterator, error) {
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
	return it.getChildIterator(index)
}

func (it *RiscVAstIterator) getChildIterator(index uint) (*RiscVAstIterator, error) {
	if index < uint(len(it.node.children)) {
		iter := makeRiscVAstIterator(it.node.children[index])
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
	nodeKind  Assembler.AstNodeKind
	children  []*RiscVAstNode
}

func makeRiscVAstNode(parent *RiscVAstNode, lineCount Assembler.LineCount, data Token, kind Assembler.AstNodeKind) RiscVAstNode {
	ast := RiscVAstNode{
		parent:    parent,
		lineCount: lineCount,
		data:      data,
		children:  []*RiscVAstNode{},
	}
	return ast
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

/*GetNodeKind returns the kind of the node within the Ast.
The Node may be a Token node, in which case the internal Token is the data.
However, if the Node is not a Token, then the internal Token is invalid*/
func (node *RiscVAstNode) GetNodeKind() Assembler.AstNodeKind {
	return node.nodeKind
}
