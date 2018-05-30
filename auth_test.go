package vault

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppRoleLogin(t *testing.T) {

	// echoHandler, passes back form parameter p
	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		res := `
{
  "auth": {
    "renewable": true,
    "lease_duration": 1200,
    "metadata": null,
    "policies": [
      "default"
    ],
    "accessor": "fd6c9a00-d2dc-3b11-0be5-af7ae0e1d374",
    "client_token": "5b1a0318-679c-9c45-e5c6-d1b9a9035d49"
  },
  "warnings": null,
  "wrap_info": null,
  "data": null,
  "lease_duration": 0,
  "renewable": false,
  "lease_id": ""
}`

		fmt.Fprint(w, res)
	}

	// create test server with handler
	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	// new client
	client := NewClient(ts.URL, true)
	err := client.AppRoleLogin("testroleid", "testsecretid")

	if err != nil {
		t.Fatalf("Failed to login: %s", err)
	}
}
