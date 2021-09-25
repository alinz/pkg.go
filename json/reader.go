package json

import (
	"encoding/json"
	"io"
)

func Reader(ptr interface{}) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		err := json.NewEncoder(pw).Encode(ptr)
		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	return pr
}
