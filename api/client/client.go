// Copyright 2023 BeyondTrust. All rights reserved.
// Package client implements functions to call Beyondtrust Secret Safe API.

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"terraform-provider-passwordsafe/api/client/entities"

	"software.sslmate.com/src/go-pkcs12"
)

var signInCount uint64
var mu sync.Mutex
var mu_out sync.Mutex

type Client struct {
	url            string
	apiKey         string
	apiAccountName string
	httpClient     *http.Client
}

// NewClient returns a http.Transport object to call secret safe API.
func NewClient(url string, apiKey string, apiAccountName string, verifyca bool, clientCertificatePath string, clientCertificateName string, clientCertificatePassword string) (*Client, error) {

	var cert tls.Certificate

	if clientCertificatePath != "" {
		pfxFile, err := ioutil.ReadFile(filepath.Join(clientCertificatePath, clientCertificateName))
		if err != nil {
			return nil, err
		}

		pfxFileBlock, err := pkcs12.ToPEM(pfxFile, clientCertificatePassword)
		if err != nil {
			return nil, err
		}

		var keyBlock, certificateBlock *pem.Block
		for _, pemBlock := range pfxFileBlock {
			if pemBlock.Type == "PRIVATE KEY" {
				keyBlock = pemBlock
			} else if pemBlock.Type == "CERTIFICATE" {
				certificateBlock = pemBlock
			}
		}

		if keyBlock == nil {
			return nil, errors.New("Error getting Key Block")
		}
		if certificateBlock == nil {
			return nil, errors.New("Error getting Certificate Block")
		}

		privateKeyData := pem.EncodeToMemory(keyBlock)
		certData := pem.EncodeToMemory(certificateBlock)

		cert, err = tls.X509KeyPair([]byte(certData), []byte(privateKeyData))

		if err != nil {
			return nil, err
		}
	}

	// TSL Config
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: !verifyca,
			Certificates:       []tls.Certificate{cert},
		},
	}

	var jar, _ = cookiejar.New(nil)

	// Client
	var client = &http.Client{
		Transport: tr,
		Jar:       jar,
	}
	return &Client{
		url:            url,
		apiKey:         apiKey,
		apiAccountName: apiAccountName,
		httpClient:     client,
	}, nil
}

/******************************************* ManageAccountFlow Methods *******************************************/

// ManageAccountFlow returns value for a specific System Name and Account Name.
func (c *Client) ManageAccountFlow(systemName string, accountName string, paths map[string]string) (string, error) {

	if len(paths) == 0 {
		paths["SignAppinPath"] = "Auth/SignAppin"
		paths["SignAppOutPath"] = "Auth/Signout"
		paths["ManagedAccountGetPath"] = fmt.Sprintf("ManagedAccounts?systemName=%v&accountName=%v", systemName, accountName)
		paths["ManagedAccountCreateRequestPath"] = "Requests"
		paths["CredentialByRequestIdPath"] = "Credentials/%v"
		paths["ManagedAccountRequestCheckInPath"] = "Requests/%v/checkin"
	}

	systemName = strings.TrimSpace(systemName)
	accountName = strings.TrimSpace(accountName)
	SignAppinOutUrl := c.requestPath(paths["SignAppOutPath"])

	if systemName == "" {
		return "", errors.New("Please use a valid system_name value")
	}

	if accountName == "" {
		return "", errors.New("Please use a valid account_name value")
	}

	SignAppinUrl := c.requestPath(paths["SignAppinPath"])
	_, err := c.SignAppin(SignAppinUrl)

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	ManagedAccountGetUrl := c.requestPath(paths["ManagedAccountGetPath"])
	managedAccount, err := c.ManagedAccountGet(systemName, accountName, ManagedAccountGetUrl)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	ManagedAccountCreateRequestUrl := c.requestPath(paths["ManagedAccountCreateRequestPath"])
	requestId, err := c.ManagedAccountCreateRequest(managedAccount.SystemId, managedAccount.AccountId, ManagedAccountCreateRequestUrl)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	CredentialByRequestIdUrl := c.requestPath(fmt.Sprintf(paths["CredentialByRequestIdPath"], requestId))
	secret, err := c.CredentialByRequestId(requestId, CredentialByRequestIdUrl)
	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	ManagedAccountRequestCheckInPath := fmt.Sprintf(paths["ManagedAccountRequestCheckInPath"], requestId)
	ManagedAccountRequestCheckInUrl := c.requestPath(ManagedAccountRequestCheckInPath)
	_, err = c.ManagedAccountRequestCheckIn(requestId, ManagedAccountRequestCheckInUrl)

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}
	c.SignOut(SignAppinOutUrl)
	secretValue, _ := strconv.Unquote(secret)
	return secretValue, nil
}

func (c *Client) ManagedAccountGet(systemName string, accountName string, url string) (entities.ManagedAccount, error) {
	// log.debug("ManagedAccountGet")
	body, err := c.httpRequest(url, "GET", bytes.Buffer{})
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

// ManagedAccountCreateRequest calls Secret Safe API Requests enpoint and returns a request Id as string.
func (c *Client) ManagedAccountCreateRequest(systemName int, accountName int, url string) (string, error) {

	data := fmt.Sprintf(`{"SystemID":%v, "AccountID":%v, "DurationMinutes":5, "Reason":"Tesr", "ConflictOption": "reuse"}`, systemName, accountName)
	b := bytes.NewBufferString(data)

	body, err := c.httpRequest(url, "POST", *b)
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

// CredentialByRequestId calls Secret Safe API Credentials/<request_id>
// enpoint and returns secret value by request Id.
func (c *Client) CredentialByRequestId(requestId string, url string) (string, error) {

	body, err := c.httpRequest(url, "GET", bytes.Buffer{})
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

// ManagedAccountRequestCheckIn calls Secret Safe API "Requests/<request_id>/checkin enpoint.
func (c *Client) ManagedAccountRequestCheckIn(requestId string, url string) (string, error) {

	data := "{}"
	b := bytes.NewBufferString(data)
	_, err := c.httpRequest(url, "PUT", *b)
	if err != nil {
		return "", err
	}

	return "", nil
}

/******************************************* SecretFlow Methods *******************************************/

// SecretFlow returns secret value for a specific path and title.
func (c *Client) SecretFlow(secretPath string, secretTitle string, separator string, paths map[string]string) (string, error) {

	if len(paths) == 0 {
		paths["SignAppinPath"] = "Auth/SignAppin"
		paths["SignAppOutPath"] = "Auth/Signout"
		paths["SecretGetSecretByPathPath"] = fmt.Sprintf("secrets-safe/secrets?title=%v&path=%v&separator=%v", secretTitle, secretPath, separator)
		paths["SecretGetFileSecretPath"] = "secrets-safe/secrets/%v/file/download"
	}

	secretPath = strings.TrimSpace(secretPath)
	secretTitle = strings.TrimSpace(secretTitle)
	separator = strings.TrimSpace(separator)
	SignAppinOutUrl := c.requestPath(paths["SignAppOutPath"])

	if secretPath == "" {
		return "", errors.New("Please use a valid Path value")
	}

	if secretTitle == "" {
		return "", errors.New("Please use a valid Title value")
	}

	SignAppinUrl := c.requestPath(paths["SignAppinPath"])
	_, err := c.SignAppin(SignAppinUrl)

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	SecretGetSecretByPathUrl := c.requestPath(paths["SecretGetSecretByPathPath"])
	secret, err := c.SecretGetSecretByPath(secretPath, secretTitle, separator, SecretGetSecretByPathUrl)

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut(SignAppinOutUrl)
		return "", errors.New(error_message)
	}

	// When secret type is FILE, it calls SecretGetFileSecret method.
	if strings.ToUpper(secret.SecretType) == "FILE" {

		SecretGetFileSecretUrl := c.requestPath(fmt.Sprintf(paths["SecretGetFileSecretPath"], secret.Id))
		fileSecretContent, err := c.SecretGetFileSecret(secret.Id, SecretGetFileSecretUrl)
		if err != nil {
			error_message := err.Error()
			fmt.Println(error_message)
			c.SignOut(SignAppinOutUrl)
			return "", errors.New(error_message)
		}
		c.SignOut(SignAppinOutUrl)
		return fileSecretContent, nil
	}

	c.SignOut(SignAppinOutUrl)
	return secret.Password, nil

}

// SecretGetSecretByPath returns secret object for a specific path, title.
func (c *Client) SecretGetSecretByPath(secretPath string, secretTitle string, separator string, url string) (entities.Secret, error) {

	body, err := c.httpRequest(url, "GET", bytes.Buffer{})
	if err != nil {
		return entities.Secret{}, err
	}

	if err != nil {
		return entities.Secret{}, err
	}

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return entities.Secret{}, err
	}

	var SecretObjectList []entities.Secret
	err = json.Unmarshal([]byte(bodyBytes), &SecretObjectList)
	if err != nil {
		return entities.Secret{}, errors.New(err.Error() + ", Ensure Password Safe version is 23.1 or greater.")
	}
	return SecretObjectList[0], nil
}

// SecretGetFileSecret call secrets-safe/secrets/<secret_id>/file/download enpoint
// and returns file secret value.
func (c *Client) SecretGetFileSecret(secretId string, url string) (string, error) {

	body, err := c.httpRequest(url, "GET", bytes.Buffer{})

	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	responseData, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}

	responseString := string(responseData)
	return responseString, nil

}

/******************************************* Common Methods *******************************************/

func (c *Client) SignAppin(url string) (entities.User, error) {
	mu.Lock()
	var userObject entities.User
	if atomic.LoadUint64(&signInCount) > 0 {
		atomic.AddUint64(&signInCount, 1)
		// TODO: log debug log("Already signed in ", atomic.LoadUint64(&signInCount))
		mu.Unlock()
		return userObject, nil
	}

	body, err := c.httpRequest(url, "POST", bytes.Buffer{})
	if err != nil {
		return entities.User{}, err
	}

	atomic.AddUint64(&signInCount, 1)
	// TODO: log debug log("signin user ", atomic.LoadUint64(&signInCount)))
	mu.Unlock()

	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return entities.User{}, err
	}

	json.Unmarshal(bodyBytes, &userObject)
	return userObject, nil
}

// SignOut signs out Secret Safe API.
// Warn: should only be called one time for all data sources.
func (c *Client) SignOut(url string) error {
	mu_out.Lock()
	if atomic.LoadUint64(&signInCount) > 1 {
		// TODO: log debug log("Ignore signout ", atomic.LoadUint64(&signInCount)))
		// decrement counter, don't signout.
		atomic.AddUint64(&signInCount, ^uint64(0))
		mu_out.Unlock()
		return nil
	}

	fmt.Println(url)
	_, err := c.httpRequest(url, "POST", bytes.Buffer{})
	if err != nil {
		return err
	}

	// decrement counter
	// TODO: log debug log("signout user ", atomic.LoadUint64(&signInCount)))
	atomic.AddUint64(&signInCount, ^uint64(0))
	mu_out.Unlock()
	return nil
}

// httpRequest template for Secret Safe API requests.
func (c *Client) httpRequest(url string, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}

	var authorizationHeader string = fmt.Sprintf("PS-Auth key=%v;runas=%v;", c.apiKey, c.apiAccountName)

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

// requestPath Build endpint path.
func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%v/%v", c.url, path)
}
