package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"terraform-provider-passwordsafe/api/client/entities"
)

type Client struct {
	url         string
	apitoken    string
	accountname string
	httpClient  *http.Client
}

type User struct {
	UserId       int    `json:"UserId"`
	EmailAddress string `json:"EmailAddress"`
	UserName     string `json:"UserName"`
	Name         string `json:"Name"`
}

func NewClient(url string, apitoken string, accountname string) *Client {

	// TSL Config
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation: tls.RenegotiateOnceAsClient,
		},
	}

	var jar, _ = cookiejar.New(nil)

	// Client
	var client = &http.Client{
		Transport: tr,
		Jar:       jar,
	}
	return &Client{
		url:         url,
		apitoken:    apitoken,
		accountname: accountname,
		httpClient:  client,
	}
}

func (c *Client) ManageAccountFlow(systemName string, accountName string) (string, error) {

	systemName = strings.TrimSpace(systemName)
	accountName = strings.TrimSpace(accountName)

	if systemName == "" {
		return "", errors.New("Please use a valid system_name value")
	}

	if accountName == "" {
		return "", errors.New("Please use a valid account_name value")
	}

	_, err := c.SignAppin()

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	managedAccount, err := c.ManagedAccountGet(systemName, accountName)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	requestId, err := c.ManagedAccountCreateRequest(managedAccount.SystemId, managedAccount.AccountId)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	secret, err := c.ManagedAccountGetCredentialByRequestId(requestId)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	_, err = c.ManagedAccountRequestCheckIn(requestId)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}
	secretValue, _ := strconv.Unquote(secret)
	return secretValue, nil
}

func (c *Client) SignAppin() (User, error) {
	body, err := c.httpRequest("Auth/SignAppin", "POST", bytes.Buffer{})
	if err != nil {
		return User{}, err
	}
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return User{}, err
	}

	var userObject User
	json.Unmarshal(bodyBytes, &userObject)
	return userObject, nil
}

func (c *Client) SignOut() error {
	_, err := c.httpRequest("Auth/SignAppin", "POST", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ManagedAccountGet(systemName string, accountName string) (entities.ManagedAccount, error) {

	path := fmt.Sprintf("ManagedAccounts?systemName=%v&accountName=%v", systemName, accountName)

	body, err := c.httpRequest(path, "GET", bytes.Buffer{})
	if err != nil {
		return entities.ManagedAccount{}, err
	}

	if err != nil {
		return entities.ManagedAccount{}, err
	}

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return entities.ManagedAccount{}, err
	}

	var managedAccountObject entities.ManagedAccount
	json.Unmarshal(bodyBytes, &managedAccountObject)
	return managedAccountObject, nil

}

func (c *Client) ManagedAccountCreateRequest(systemName int, accountName int) (string, error) {

	data := fmt.Sprintf(`{"SystemID":%v, "AccountID":%v, "DurationMinutes":5, "Reason":"Tesr", "ConflictOption": "reuse"}`, systemName, accountName)
	b := bytes.NewBufferString(data)

	path := fmt.Sprintf("Requests")

	body, err := c.httpRequest(path, "POST", *b)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return "", err
	}

	responseString := string(bodyBytes)

	return responseString, nil

}

func (c *Client) ManagedAccountGetCredentialByRequestId(requestId string) (string, error) {

	path := fmt.Sprintf("Credentials/%v", requestId)

	body, err := c.httpRequest(path, "GET", bytes.Buffer{})
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return "", err
	}

	responseString := string(bodyBytes)

	return responseString, nil

}

func (c *Client) ManagedAccountRequestCheckIn(requestId string) (string, error) {

	data := "{}"
	b := bytes.NewBufferString(data)

	path := fmt.Sprintf("Requests/%v/checkin", requestId)

	_, err := c.httpRequest(path, "PUT", *b)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) httpRequest(path string, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	url := c.requestPath(path)

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}

	var authorizationHeader string = fmt.Sprintf("PS-Auth key=%v;runas=%v;", c.apitoken, c.accountname)

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authorizationHeader},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)

		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}

	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%v/%v", c.url, path)
}
