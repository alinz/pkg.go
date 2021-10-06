package httputil_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/alinz/pkg.go/httputil"
	"github.com/stretchr/testify/assert"
)

func TestUnixSocket(t *testing.T) {
	server := httputil.NewServer("unix", "/tmp/sample.sock", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	go func() {
		err := server.Start(context.Background())
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	client := httputil.NewClient("unix", "/tmp/sample.sock")
	resp, err := client.Get("http://localhost")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	server.Stop(context.Background())
	time.Sleep(1 * time.Second)
}
