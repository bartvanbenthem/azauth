package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Requester interface {
	GetToken() (Response, error)
}

type Response struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

type Credential struct {
	TenantID      string
	ApplicationID string
	ClientSecret  string
}

type RMClient struct {
	Auth     Credential
	Response Response
}

// RMClient.GetToken gets and returns an access token from the Microsoft Resource Manager API
func (c *RMClient) GetToken() (Response, error) {
	resp := c.Response
	resp, err := tokenRequest(c.Auth, "https://management.azure.com")
	if err != nil {
		return resp, err
	}
	return resp, err
}

type GraphClient struct {
	Auth     Credential
	Response Response
}

// GraphClient.GetToken gets and returns an access token from the Microsoft Graph API
func (c *GraphClient) GetToken() (Response, error) {
	resp := c.Response
	resp, err := tokenRequest(c.Auth, "https://graph.microsoft.com")
	if err != nil {
		return resp, err
	}
	return resp, err
}

// tokenRequest is a generic function that based on resource input returns the requested token
// used by the GetToken methods for different resources
func tokenRequest(c Credential, resource string) (Response, error) {
	token := Response{}
	path := fmt.Sprintf("/%v/oauth2/token", c.TenantID)
	data := url.Values{}

	data.Add("grant_type", "client_credentials")
	data.Add("client_id", c.ApplicationID)
	data.Add("client_secret", c.ClientSecret)
	data.Add("resource", resource)

	u, err := url.ParseRequestURI("https://login.microsoftonline.com")
	if err != nil {
		return token, err
	}

	u.Path = path
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	httpClient := &http.Client{Timeout: time.Second * 10}
	resp, err := httpClient.Do(req)
	if err != nil {
		return token, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	json.Unmarshal(body, &token)
	return token, err
}
