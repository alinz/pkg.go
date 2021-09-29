package input

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func FromUser(message string, secure bool) (string, error) {
	fmt.Print(message)

	var value []byte
	var err error

	if secure {
		value, err = term.ReadPassword(0)
		fmt.Println("")
	} else {
		reader := bufio.NewReader(os.Stdin)
		value, _, err = reader.ReadLine()
	}

	return string(value), err
}
