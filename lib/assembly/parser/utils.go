package parser

import Assembler "github.com/chenhowa/computer/lib/assembly"

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
	return tokenType == Assembler.Label
}

func isMnemonic(token Token) bool {
	tokenType := token.GetTokenType()

	return tokenType <= Assembler.BLEU && tokenType >= Assembler.RDCYCLE
}
