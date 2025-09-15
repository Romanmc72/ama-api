package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: Get these to be dynamic per-test-run
const (
	UserEmail          = "t@t.com"
	UserPass           = "password123"
	EmulatorHost       = "localhost"
	EmulatorPort       = 9099
	ResourceServerPort = 8088
	ResourceServerHost = "localhost"
	ApiKey             = "fake-api-key"
	IsSecure           = false
)

type ReturnedToken struct {
	Kind         string `json:"kind"`
	IsNewUser    bool   `json:"isNewUser"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpirationIn string `json:"expiresIn"`
}

func HitApi[T any](client *http.Client, url string, method string, token string, payload any, responseObj T) (T, error) {
	var req *http.Request
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return responseObj, err
		}
		req, err = http.NewRequest(method, url, strings.NewReader(string(body)))
		if err != nil {
			return responseObj, err
		}
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return responseObj, err
		}
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	resp, err := client.Do(req)
	if err != nil {
		return responseObj, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return responseObj, fmt.Errorf("unsuccessful response code: %d for method: %s url: %s", resp.StatusCode, method, url)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseObj, err
	}
	err = json.Unmarshal(respBody, &responseObj)
	if err != nil {
		return responseObj, err
	}
	return responseObj, nil
}
