package govault

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	Address    string
	Token      string
	HttpClient *http.Client
}

func NewClient(addr string, tlsSkipVerify bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsSkipVerify},
	}
	client := &http.Client{Transport: tr}

	return &Client{
		Address:    addr,
		HttpClient: client,
	}
}
