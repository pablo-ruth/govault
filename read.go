package govault

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Read(path string, code int) (map[string]interface{}, error) {
	if c.token == "" {
		return nil, fmt.Errorf("empty vault token")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLSSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/v1/%s", c.Address, path),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Vault-Token", c.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != code {
		return nil, fmt.Errorf("read failed with http code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var read map[string]interface{}
	if err = json.Unmarshal(body, &read); err != nil {
		return nil, err
	}

	return read["data"].(map[string]interface{}), nil
}
