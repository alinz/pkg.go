package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/alinz/pkg.go/token"
)

//Lexer lexer struct
type Lexer struct {
	input   string
	start   int
	pos     int
	width   int
	emitter token.Emitter
}

func (l *Lexer) Emit(tokenType token.Type) {
	token := &token.Token{
		Type:  tokenType,
		Value: l.input[l.start:l.pos],
		Start: l.start,
		End:   l.pos,
	}
	l.emitter.Emit(token)
	l.start = l.pos
}

func (l *Lexer) Next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// PeekN is a function that returns the next runes in the input
// based on given number without advancing the position
// returns two values, the runes and the number of iteration <= n
func (l *Lexer) PeekN(n int) (string, int) {
	sb := strings.Builder{}

	total := 0
	i := 0

	for i < n {
		value := l.Next()
		if value == 0 {
			break
		}

		total += l.width
		sb.WriteRune(value)
		i++
	}

	l.pos -= total

	return sb.String(), i
}

func (l *Lexer) Backup() {
	l.pos -= l.width
	l.width = 0
}

func (l *Lexer) Ignore() {
	l.start = l.pos
}

func (l *Lexer) Accept(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

func (l *Lexer) AcceptRun(valid string) {
	for strings.ContainsRune(valid, l.Next()) {
	}
	l.Backup()
}

func (l *Lexer) AcceptRunUntil(invalid string) {
	for {
		next := l.Next()
		if next == 0 || strings.ContainsRune(invalid, next) {
			break
		}
	}
	l.Backup()
}

func (l *Lexer) Errorf(format string, args ...interface{}) {
	l.emitter.Emit(&token.Token{
		Type:  token.Error,
		Name:  "Error",
		Value: fmt.Sprintf(format, args...),
		Start: l.start,
		End:   l.pos,
	})
}

func (l *Lexer) Run(state State) {
	for state != nil {
		state = state(l)
	}
}

type State func(*Lexer) State

func New(input string, emitter token.Emitter) *Lexer {
	return &Lexer{
		input:   input,
		emitter: emitter,
	}
}
