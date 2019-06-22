package parser

import (
	"fmt"
	"strings"

	Assembler "github.com/chenhowa/computer/lib/assembly"
)

func isOperand(token Token) bool {
	tokenType := token.GetTokenType()

	return isRegister(tokenType) ||
		isNumericConstant(tokenType) ||
		isRegisterAndImmediate(tokenType) ||
		isIdentifier(tokenType)
}

func isRegister(tokenType Assembler.TokenType) bool {
	return tokenType <= Assembler.X31 && tokenType >= Assembler.X0
}

func isNumericConstant(tokenType Assembler.TokenType) bool {
	return tokenType == Assembler.NumericConstant
}

func isRegisterAndImmediate(tokenType Assembler.TokenType) bool {
	return tokenType == Assembler.RegisterAndImmediate
}

func isIdentifier(tokenType Assembler.TokenType) bool {
	return tokenType == Assembler.Identifier
}

func isLabel(token Token) bool {
	return token.GetTokenType() == Assembler.Label
}

func isMnemonic(token Token) bool {
	tokenType := token.GetTokenType()

	return tokenType <= Assembler.BLEU && tokenType >= Assembler.RDCYCLE
}

func mnemonic(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	token, tokenErr := stream.Next()

	if tokenErr == nil && isMnemonic(token) {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		reset.Reset()
		ast := RiscVAst{}
		return ast, false
	}
}

func operand(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	token, tokenErr := stream.Next()

	if tokenErr == nil && isOperand(token) {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		reset.Reset()
		ast := RiscVAst{}
		return ast, false
	}
}

func optionalNewlines(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	optionalNewlinesAst, ok := _optionalNewlines(stream)

	// Since this is optional newlines, as long as there was no stream error,
	// it is considered a successful parse by the CALLER. The return values have nothing to do with this.
	// But how would you even know if there was a stream error? You couldn't.
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

func newline(stream tokenStream) (tree RiscVAst, success bool) {
	reset := stream.Save()
	token, tokenErr := stream.Next()
	if tokenErr == nil && token.GetTokenType() == Assembler.Newline {
		node := makeRiscVAstNode(nil, 1, token, Assembler.Token)
		ast := RiscVAst{
			root: &node,
		}
		return ast, true
	} else {
		reset.Reset()
		ast := RiscVAst{}
		return ast, false
	}
}

func convertToString(iter AstIterator) string {
	var builder strings.Builder
	builder.WriteString("(" + getIterRepr(iter))
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

func getIterRepr(iter AstIterator) (str string) {
	defer func() {
		if r := recover(); r != nil {
			node := iter.GetAstNode()
			str = fmt.Sprintf("%d", uint(node.GetNodeKind()))
		}
	}()

	return iter.GetAstNode().GetTokenString()
}
