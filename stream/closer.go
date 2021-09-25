package stream

import "io"

type readCloser struct {
	r io.Reader
	c io.Closer
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	return rc.r.Read(p)
}

func (rc *readCloser) Close() error {
	return rc.c.Close()
}

func NewCloser(r io.Reader, c io.Closer) io.ReadCloser {
	return &readCloser{r, c}
}
