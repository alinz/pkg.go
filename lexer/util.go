package lexer

import (
	"unicode"
)

func isAlphaNumeric(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

func ignoreWhiteSpace(l *Lexer) {
	l.AcceptRun(` \t\n`)
	l.Ignore()
}
