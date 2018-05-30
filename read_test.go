package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRead(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		res := `
{
  "auth": null,
  "data": {
    "foo": "bar"
  },
  "lease_duration": 2764800,
  "lease_id": "",
  "renewable": false
}
`

		fmt.Fprint(w, res)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(testHandler))
	defer ts.Close()

	client := NewClient(ts.URL, false)
	client.Token = "testtoken"

	entries, err := client.Read("secrets", 200)
	if err != nil {
		t.Fatalf("Failed to read: %s", err)
	}

	if entries["foo"] != "bar" {
		t.Fatal("Failed to check read foo entry")
	}
}
