package stream

import (
	"bytes"
	"io"
)

// Peek n bytes from the reader, you could use bufio.NewReader to peek sa well
func Peek(r io.Reader, n int64) ([]byte, io.Reader) {
	if n <= 0 {
		return []byte{}, r
	}

	buffer := bytes.Buffer{}
	_, err := io.Copy(&buffer, io.LimitReader(r, int64(n)))
	if err != nil {
		r = Error(err)
	}

	data := buffer.Bytes()
	return data, io.MultiReader(bytes.NewReader(data), r)
}

// PeekCloser is a io.ReadCloser that peeks n bytes from the reader
func PeekCloser(rc io.ReadCloser, n int64) ([]byte, io.ReadCloser) {
	peek, r := Peek(rc, n)
	return peek, NewCloser(r, rc)
}
