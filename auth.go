package govault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) AppRoleLogin(roleid, secretid string) error {
	if c.Address == "" {
		return fmt.Errorf("Empty vault address")
	}
	if roleid == "" {
		return fmt.Errorf("Empty roleid")
	}
	if secretid == "" {
		return fmt.Errorf("Empty secretid")
	}

	data, err := json.Marshal(map[string]string{
		"role_id":   roleid,
		"secret_id": secretid,
	})
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLSSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/auth/approle/login", c.Address),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("approle login failed with http code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResult map[string]interface{}
	if err := json.Unmarshal(body, &loginResult); err != nil {
		return err
	}

	authRes, ok := loginResult["auth"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON")
	}

	clientToken, ok := authRes["client_token"].(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON")
	}

	c.token = clientToken
	return nil
}

func (c *Client) Logout() {
	if c.token == "" {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.TLSSkipVerify},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/auth/token/revoke-self", c.Address),
		nil,
	)
	if err != nil {
		return
	}

	req.Header.Add("X-Vault-Token", c.token)
	_, err = client.Do(req)
	if err != nil {
		log.Print("Logout failed")
	}
}
