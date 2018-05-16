package govault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListKeys(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		res := `
{
  "auth": null,
  "data": {
    "keys": ["first", "foo/", "bar"]
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

	list, err := client.ListKeys("secrets", 200)
	if err != nil {
		t.Fatalf("Failed to get list of keys: %s", err)
	}

	if strings.Join(list, ",") != "first,bar" {
		t.Fatal("Wrong list of keys", strings.Join(list, ","))
	}
}
