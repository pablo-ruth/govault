package govault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Write(path string, data map[string]interface{}, code int) error {
	if c.token == "" {
		return fmt.Errorf("empty vault token")
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLSSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/%s", c.Address, path),
		bytes.NewBuffer(dataJSON),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", c.token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != code {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("write failed on %s with http code %d: %s", path, resp.StatusCode, body)
	}

	return nil
}
