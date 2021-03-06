package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	Assembler "github.com/chenhowa/computer/lib/assembly"
)

func isNumericConstant(tokenString string) bool {
	var rg = regexp.MustCompile(`(^0$)|(^[1-9]\d*$)|(^[1-9](\d|\d\d)?(,\d\d\d)*$)`)

	match := rg.MatchString(tokenString)
	return match
}

func isLabel(tokenString string) bool {
	var rg = regexp.MustCompile(`^[A-Z][a-z]*:$`)
	return rg.MatchString(tokenString)
}

func isUnskippableChar(val byte) bool {
	if unicode.IsLetter(rune(val)) {
		return true
	}

	if unicode.IsNumber(rune(val)) {
		return true
	}

	if val == '\n' {
		return true
	}

	if val == ',' {
		return true
	}

	if val == ':' {
		return true
	}

	if val == '(' || val == ')' {
		return true
	}

	return false
}

func getCleanRiscVToken(tokenType Assembler.TokenType, tokenString string, charsSinceLastNewline Assembler.CharCount) RiscVToken {
	token := makeRiscVToken(tokenType, cleanTokenString(tokenType, tokenString), charsSinceLastNewline)
	return token
}

func cleanTokenString(tokenType Assembler.TokenType, tokenString string) string {
	if tokenType == Assembler.NumericConstant {
		return strings.Replace(tokenString, ",", "", -1)
	}

	if tokenType == Assembler.Label {
		return strings.Replace(tokenString, ":", "", -1)
	}

	return tokenString
}

func continueReadingTokenInput(latestChar byte, readInput string) bool {
	return !suddenNewline(readInput, latestChar) && isUnskippableChar(latestChar)
}

func suddenNewline(readInput string, latestChar byte) bool {
	return suddenChar('\n', readInput, latestChar)
}

func suddenComma(readInput string, latestChar byte) bool {
	return suddenChar(',', readInput, latestChar)
}

func suddenChar(char byte, readInput string, latestChar byte) bool {
	return (uint(len(readInput)) > 0) && (latestChar == char)
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

	tokenType, ok = registerToToken[Register(tokenString)]
	if isRegister(tokenString) && ok {
		return tokenType, nil
	}

	if isRegisterImmediate(tokenString) {
		return Assembler.RegisterAndImmediate, nil
	}

	tokenType, ok = registerNicknameToToken[RegisterNickname(tokenString)]
	if ok {
		return tokenType, nil
	}

	// This check should go last, as most things will register as an identifier
	if isIdentifier(tokenString) {
		return Assembler.Identifier, nil
	}

	return tokenType, fmt.Errorf("getTokenType: no token type found for this token %s", tokenString)
}

func isRegister(tokenString string) bool {
	var rg = regexp.MustCompile(`^x(\d|([1-3]\d))$`)

	match := rg.MatchString(tokenString)
	return match
}

func isRegisterImmediate(tokenString string) bool {
	numericConstant := `((0)|([1-9]\d*)|([1-9](\d|\d\d)?(,\d\d\d)*))`
	register := `(\(x(\d|([1-2]\d)|(3[0-1]))\))`

	var rg = regexp.MustCompile(`^` + numericConstant + register + `$`)
	return rg.MatchString(tokenString)
}

func isIdentifier(tokenString string) bool {
	var rg = regexp.MustCompile(`^[A-Za-z]+$`)
	return rg.MatchString(tokenString)
}
