package govault

import "testing"

func TestNewClient(t *testing.T) {
	addr := "https://myvaultsrv.test"
	tlsSkipVerify := true

	client := NewClient(addr, tlsSkipVerify)

	if client.Address != addr {
		t.Fatal("Failed to set address on new client")
	}
	if client.HttpClient == nil {
		t.Fatal("Failed to create httpclient on new client")
	}
}
