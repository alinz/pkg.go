package httputil_test

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/alinz/pkg.go/httputil"
	"github.com/stretchr/testify/assert"
)

func TestUnixSocket(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "sample.sock")

	server := httputil.NewServer("unix", sockPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Start(ctx)
	assert.NoError(t, err)

	client := httputil.NewClient("unix", sockPath)
	resp, err := client.Get("http://localhost")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = server.Stop(context.Background())
	assert.NoError(t, err)
}
