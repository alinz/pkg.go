package json

import (
	"encoding/json"
	"io"
)

func ParseReader(r io.Reader, ptr interface{}) error {
	if rd, ok := r.(io.ReadCloser); ok {
		defer rd.Close()
	}
	return json.NewDecoder(r).Decode(ptr)
}

func Writer(w io.Writer, ptr interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // need this to disable encoding & to unicode
	return encoder.Encode(ptr)
}
