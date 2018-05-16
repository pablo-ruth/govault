package govault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		res := ""
		fmt.Fprint(w, res)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	client := NewClient(ts.URL, false)
	client.Token = "testtoken"

	err := client.Write("secrets/mysecret", map[string]string{"myuser": "testuser"}, 200)
	if err != nil {
		t.Fatalf("Failed to write: %s", err)
	}
}
