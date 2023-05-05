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
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"terraform-provider-passwordsafe/api/client/entities"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/crypto/pkcs12"
)

var signInCount uint64

var logger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)

type Client struct {
	url                string
	apiKey             string
	apiAccountName     string
	httpClient         *http.Client
	exponentialBackOff *backoff.ExponentialBackOff
	testMode           bool
}

// NewClient returns a http.Transport object to call secret safe API.
func NewClient(url string, apiKey string, apiAccountName string, verifyca bool, clientCertificatePath string, clientCertificateName string, clientCertificatePassword string) (*Client, error) {

	logFile, _ := os.Create("ProviderLogs.log")
	logger.SetOutput(logFile)

	var cert tls.Certificate

	if clientCertificatePath != "" {
		pfxFile, err := ioutil.ReadFile(filepath.Join(clientCertificatePath, clientCertificateName))
		if err != nil {
			Logging("ERROR", err.Error())
			return nil, err
		}

		pfxFileBlock, err := pkcs12.ToPEM(pfxFile, clientCertificatePassword)
		if err != nil {
			Logging("ERROR", err.Error())
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
			err = errors.New("Error getting Key Block")
			Logging("ERROR", err.Error())
			return nil, err
		}
		if certificateBlock == nil {
			err = errors.New("Error getting Certificate Block")
			Logging("ERROR", err.Error())
			return nil, err
		}

		privateKeyData := pem.EncodeToMemory(keyBlock)
		certData := pem.EncodeToMemory(certificateBlock)

		cert, err = tls.X509KeyPair([]byte(certData), []byte(privateKeyData))

		if err != nil {
			Logging("ERROR", err.Error())
			return nil, err
		}
	}

	// TSL Config
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: !verifyca,
			Certificates:       []tls.Certificate{cert},
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		},
	}

	var jar, _ = cookiejar.New(nil)

	// Client
	var client = &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	backoffDefinition := backoff.NewExponentialBackOff()

	testMode := false

	//Checking TEST_MODE env variable, if value is "true" means is for unit tests.
	if strings.ToLower(os.Getenv("TEST_MODE")) == "true" {
		// Test mode true for unit tests
		testMode = true
		// Configuring ExponentialBackOff object with just one retry for unit tests.
		backoffDefinition.MaxElapsedTime = time.Second
	} else {
		// Configuring ExponentialBackOff object with custom configuration for real scenario
		backoffDefinition := backoff.NewExponentialBackOff()
		backoffDefinition.InitialInterval = 1 * time.Second
		backoffDefinition.MaxElapsedTime = 15 * time.Second
		backoffDefinition.RandomizationFactor = 0.5
	}

	return &Client{
		url:                url,
		apiKey:             apiKey,
		apiAccountName:     apiAccountName,
		httpClient:         client,
		exponentialBackOff: backoffDefinition,
		testMode:           testMode,
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

	var err error

	if systemName == "" {
		err = errors.New("Please use a valid system_name value")
		Logging("ERROR", err.Error())
		return "", err
	}

	if accountName == "" {
		err = errors.New("Please use a valid system_name value")
		Logging("ERROR", err.Error())
		return "", err
	}

	ManagedAccountGetUrl := c.RequestPath(paths["ManagedAccountGetPath"])
	managedAccount, err := c.ManagedAccountGet(systemName, accountName, ManagedAccountGetUrl)
	if err != nil {
		Logging("ERROR", err.Error())
		return "", err
	}

	ManagedAccountCreateRequestUrl := c.RequestPath(paths["ManagedAccountCreateRequestPath"])
	requestId, err := c.ManagedAccountCreateRequest(managedAccount.SystemId, managedAccount.AccountId, ManagedAccountCreateRequestUrl)
	if err != nil {
		Logging("ERROR", err.Error())
		return "", err
	}

	CredentialByRequestIdUrl := c.RequestPath(fmt.Sprintf(paths["CredentialByRequestIdPath"], requestId))
	secret, err := c.CredentialByRequestId(requestId, CredentialByRequestIdUrl)
	if err != nil {
		Logging("ERROR", err.Error())
		return "", err
	}

	ManagedAccountRequestCheckInPath := fmt.Sprintf(paths["ManagedAccountRequestCheckInPath"], requestId)
	ManagedAccountRequestCheckInUrl := c.RequestPath(ManagedAccountRequestCheckInPath)
	_, err = c.ManagedAccountRequestCheckIn(requestId, ManagedAccountRequestCheckInUrl)

	if err != nil {
		Logging("ERROR", err.Error())
		return "", err
	}

	secretValue, _ := strconv.Unquote(secret)
	return secretValue, nil
}

func (c *Client) ManagedAccountGet(systemName string, accountName string, url string) (entities.ManagedAccount, error) {
	messageLog := fmt.Sprintf("%v %v", "GET", url)
	Logging("DEBUG", messageLog)

	var body io.ReadCloser
	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		body, technicalError, businessError, _ = c.callSecretSafeAPI(url, "GET", bytes.Buffer{}, "ManagedAccountGet")
		if technicalError != nil {
			return technicalError
		}
		return nil

	}, c.exponentialBackOff)

	if technicalError != nil {
		return entities.ManagedAccount{}, technicalError
	}

	if businessError != nil {
		return entities.ManagedAccount{}, businessError
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
	messageLog := fmt.Sprintf("%v %v", "POST", url)
	Logging("DEBUG", messageLog)

	data := fmt.Sprintf(`{"SystemID":%v, "AccountID":%v, "DurationMinutes":5, "Reason":"Tesr", "ConflictOption": "reuse"}`, systemName, accountName)
	b := bytes.NewBufferString(data)

	var body io.ReadCloser
	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		body, technicalError, businessError, _ = c.callSecretSafeAPI(url, "POST", *b, "ManagedAccountCreateRequest")
		return technicalError
	}, c.exponentialBackOff)

	if technicalError != nil {
		return "", technicalError
	}

	if businessError != nil {
		return "", businessError
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
	messageLog := fmt.Sprintf("%v %v", "GET", url)
	Logging("DEBUG", messageLog)

	var body io.ReadCloser
	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		body, technicalError, businessError, _ = c.callSecretSafeAPI(url, "GET", bytes.Buffer{}, "CredentialByRequestId")
		return technicalError
	}, c.exponentialBackOff)

	if technicalError != nil {
		return "", technicalError
	}

	if businessError != nil {
		return "", businessError
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
	messageLog := fmt.Sprintf("%v %v", "PUT", url)
	Logging("DEBUG", messageLog)

	data := "{}"
	b := bytes.NewBufferString(data)

	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		_, technicalError, businessError, _ = c.callSecretSafeAPI(url, "PUT", *b, "ManagedAccountRequestCheckIn")
		return technicalError
	}, c.exponentialBackOff)

	if technicalError != nil {
		return "", technicalError
	}

	if businessError != nil {
		return "", businessError
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

	var err error
	if secretPath == "" {
		err = errors.New("Please use a valid Path value")
		Logging("ERROR", err.Error())
		return "", err
	}

	if secretTitle == "" {
		err = errors.New("Please use a valid Title value")
		Logging("ERROR", err.Error())
		return "", err
	}

	SecretGetSecretByPathUrl := c.RequestPath(paths["SecretGetSecretByPathPath"])
	secret, err := c.SecretGetSecretByPath(secretPath, secretTitle, separator, SecretGetSecretByPathUrl)

	if err != nil {
		Logging("ERROR", err.Error())
		return "", err
	}

	// When secret type is FILE, it calls SecretGetFileSecret method.
	if strings.ToUpper(secret.SecretType) == "FILE" {

		SecretGetFileSecretUrl := c.RequestPath(fmt.Sprintf(paths["SecretGetFileSecretPath"], secret.Id))
		fileSecretContent, err := c.SecretGetFileSecret(secret.Id, SecretGetFileSecretUrl)
		if err != nil {
			Logging("ERROR", err.Error())
			return "", err
		}

		return fileSecretContent, nil
	}

	return secret.Password, nil

}

// SecretGetSecretByPath returns secret object for a specific path, title.
func (c *Client) SecretGetSecretByPath(secretPath string, secretTitle string, separator string, url string) (entities.Secret, error) {
	messageLog := fmt.Sprintf("%v %v", "GET", url)
	Logging("DEBUG", messageLog)

	var body io.ReadCloser
	var technicalError error
	var businessError error
	var scode int

	technicalError = backoff.Retry(func() error {
		body, technicalError, businessError, scode = c.callSecretSafeAPI(url, "GET", bytes.Buffer{}, "SecretGetSecretByPath")
		return technicalError
	}, c.exponentialBackOff)

	if technicalError != nil {
		return entities.Secret{}, technicalError
	}

	if businessError != nil {
		return entities.Secret{}, businessError
	}

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return entities.Secret{}, err
	}

	var SecretObjectList []entities.Secret
	err = json.Unmarshal([]byte(bodyBytes), &SecretObjectList)
	if err != nil {
		err = errors.New(err.Error() + ", Ensure Password Safe version is 23.1 or greater.")
		return entities.Secret{}, err
	}

	if len(SecretObjectList) == 0 {
		scode = 404
		err = fmt.Errorf("Error %v: StatusCode: %v ", "SecretGetSecretByPath, Secret was not found", scode)
		return entities.Secret{}, err
	}

	return SecretObjectList[0], nil
}

// SecretGetFileSecret call secrets-safe/secrets/<secret_id>/file/download enpoint
// and returns file secret value.
func (c *Client) SecretGetFileSecret(secretId string, url string) (string, error) {
	messageLog := fmt.Sprintf("%v %v", "GET", url)
	Logging("DEBUG", messageLog)

	var body io.ReadCloser
	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		body, technicalError, businessError, _ = c.callSecretSafeAPI(url, "GET", bytes.Buffer{}, "SecretGetFileSecret")
		return technicalError
	}, c.exponentialBackOff)

	if technicalError != nil {
		return "", technicalError
	}

	if businessError != nil {
		return "", businessError
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

	var userObject entities.User
	var body io.ReadCloser
	var technicalError error
	var businessError error
	var scode int

	err := backoff.Retry(func() error {
		body, technicalError, businessError, scode = c.callSecretSafeAPI(url, "POST", bytes.Buffer{}, "SignAppin")
		if scode == 0 {
			return nil
		}
		return technicalError
	}, c.exponentialBackOff)

	if err != nil {
		if !c.testMode {

		}
		return entities.User{}, err
	}

	if scode == 0 {
		if !c.testMode {

		}
		return entities.User{}, technicalError
	}

	if businessError != nil {
		if !c.testMode {

		}
		return entities.User{}, businessError
	}

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

	Logging("DEBUG", url)

	var technicalError error
	var businessError error

	technicalError = backoff.Retry(func() error {
		_, technicalError, businessError, _ = c.callSecretSafeAPI(url, "POST", bytes.Buffer{}, "SignOut")
		return technicalError
	}, c.exponentialBackOff)

	if businessError != nil {
		Logging("ERROR", businessError.Error())
		return businessError
	}

	return nil
}

// httpRequest template for Secret Safe API requests.
func (c *Client) httpRequest(url string, method string, body bytes.Buffer) (closer io.ReadCloser, technicalError error, businessError error, scode int) {

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err, nil, 0
	}

	var authorizationHeader string = fmt.Sprintf("PS-Auth key=%v;runas=%v;", c.apiKey, c.apiAccountName)

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authorizationHeader},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		Logging("ERROR", err.Error())
		return nil, err, nil, 0
	}

	if resp.StatusCode >= http.StatusInternalServerError || resp.StatusCode == http.StatusRequestTimeout {
		err = fmt.Errorf("Error %v: StatusCode: %v, %v, %v", method, scode, err, body)
		Logging("ERROR", err.Error())
		return nil, err, nil, resp.StatusCode
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		err = fmt.Errorf("got a non 200 status code: %v - %v", resp.StatusCode, respBody)
		return nil, nil, err, resp.StatusCode
	}

	return resp.Body, nil, nil, resp.StatusCode
}

// requestPath Build endpint path.
func (c *Client) RequestPath(path string) string {
	return fmt.Sprintf("%v/%v", c.url, path)
}

// call httpRequest method according to parameters.
func (c *Client) callSecretSafeAPI(url string, httpMethod string, body bytes.Buffer, method string) (io.ReadCloser, error, error, int) {
	response, technicalError, businessError, scode := c.httpRequest(url, httpMethod, body)
	if technicalError != nil {
		messageLog := fmt.Sprintf("Error in %v %v \n", method, technicalError)
		Logging("ERROR", messageLog)
	}

	if businessError != nil {
		messageLog := fmt.Sprintf("Error in %v: %v \n", method, businessError)
		Logging("ERROR", messageLog)
	}
	return response, technicalError, businessError, scode
}

// log message in log file.
func Logging(prefix string, message string) {
	prefix = fmt.Sprintf("%v :", prefix)
	logger.SetPrefix(prefix)
	logger.Println(message)
}
