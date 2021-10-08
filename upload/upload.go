package upload

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/alinz/pkg.go/json"
)

const payloadKey = "payload"

// Parse accepts payload which parsed form values inside payload field
// and returns ReadCloser which refers to given file
// Note that return ReadCloser object needs to be closed
func Parse(r *http.Request, payload interface{}) (io.ReadCloser, error) {
	value := r.FormValue(payloadKey)
	if value != "" {
		err := json.ParseReader(strings.NewReader(value), payload)
		if err != nil {
			return nil, err
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CreateRequest prepares a upload request
// this is multipart content upload which let's you send payload data
func CreateRequest(url string, r io.Reader, filename string, payload interface{}) (*http.Request, error) {
	pr, pw := io.Pipe()

	mpw := multipart.NewWriter(pw)

	req, err := http.NewRequest(http.MethodPost, url, pr)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mpw.FormDataContentType())

	go func() {
		var err error
		var part io.Writer

		defer pw.Close()

		var buffer bytes.Buffer

		if payload != nil {
			if err = json.Writer(&buffer, payload); err != nil {
				pw.CloseWithError(err)
				return
			}
		}

		if err = mpw.WriteField(payloadKey, buffer.String()); err != nil {
			pw.CloseWithError(err)
			return
		}

		if part, err = mpw.CreateFormFile("file", filename); err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err = io.Copy(part, r); err != nil {
			pw.CloseWithError(err)
			return
		}

		if err = mpw.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	return req, err
}
