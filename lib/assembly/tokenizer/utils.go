package tokenizer

import (
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
	return (uint(len(readInput)) > 0) && (latestChar == '\n')
}
