package upload_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alinz/pkg.go/upload"
)

type testPayload struct {
	ID string `json:"id"`
}

func TestUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p testPayload

		fmt.Println("Header", r.Header.Get("Content-Type"))

		file, err := upload.Parse(r, &p)
		if err != nil {
			t.Fatal(err)
		}

		defer file.Close()

		b, _ := ioutil.ReadAll(file)
		fmt.Println(string(b))
		fmt.Println(p)

		w.Write([]byte("done"))
	}))

	defer server.Close()

	content := strings.NewReader("Hello World")

	fmt.Println("sending")
	req, err := upload.CreateRequest(server.URL, content, "hello.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}
