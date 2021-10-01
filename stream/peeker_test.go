package stream_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/alinz/pkg.go/stream"
	"github.com/stretchr/testify/assert"
)

func TestPeeker(t *testing.T) {
	content := []byte("hello world")

	peek, r := stream.Peek(bytes.NewReader(content), 5)

	assert.Equal(t, peek, []byte("hello"))

	b, err := io.ReadAll(r)
	assert.NoError(t, err)

	assert.Equal(t, b, []byte("hello world"))
}
