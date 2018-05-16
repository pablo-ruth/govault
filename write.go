package govault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Write(path string, data map[string]string, code int) error {
	if c.Token == "" {
		return fmt.Errorf("empty vault token")
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/%s", c.Address, path),
		bytes.NewBuffer(dataJSON),
	)
	if err != nil {
		return err
	}

	req.Header.Add("X-Vault-Token", c.Token)
	resp, err := c.HttpClient.Do(req)
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
