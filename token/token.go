package token

import "fmt"

type Type byte

type Token struct {
	Type  Type
	Name  string
	Value string
	Start int
	End   int
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s, Value: %s, Start: %d, End: %d", t.Name, t.Value, t.Start, t.End)
}

const (
	EOF = iota
	Error
)
