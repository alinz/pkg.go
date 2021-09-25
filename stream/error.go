package stream

import "io"

type errReader struct{ err error }

func (er *errReader) Read([]byte) (int, error) { return 0, er.err }

func Error(err error) io.Reader { return &errReader{err} }
