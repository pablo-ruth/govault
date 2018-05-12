package govault

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) ListKeys(path string, code int) ([]string, error) {
	if c.token == "" {
		return nil, fmt.Errorf("empty vault token")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLSSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(
		"LIST",
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
		return nil, fmt.Errorf("list failed with http code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list map[string]interface{}
	if err = json.Unmarshal(body, &list); err != nil {
		return nil, err
	}

	data := list["data"].(map[string]interface{})

	var keys []string
	for _, key := range data["keys"].([]interface{}) {
		keyStr := key.(string)
		if last := len(keyStr) - 1; last >= 0 && keyStr[last] != '/' {
			keys = append(keys, key.(string))
		}
	}

	return keys, nil
}
