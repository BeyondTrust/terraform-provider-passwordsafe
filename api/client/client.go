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
	"terraform-provider-passwordsafe/api/client/entities"

	"software.sslmate.com/src/go-pkcs12"
)

type Client struct {
	url            string
	apiKey         string
	apiAccountName string
	httpClient     *http.Client
}

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
		if err := ioutil.WriteFile(filepath.Join(clientCertificatePath, "certificate_key.pem"), privateKeyData, 0644); err != nil {
			return nil, err
		}

		certData := pem.EncodeToMemory(certificateBlock)
		if err := ioutil.WriteFile(filepath.Join(clientCertificatePath, "certificate_cer.pem"), certData, 0644); err != nil {
			return nil, err
		}
		cert, err = tls.LoadX509KeyPair(filepath.Join(clientCertificatePath, "certificate_cer.pem"), filepath.Join(clientCertificatePath, "certificate_key.pem"))
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
	c.SignOut()
	return secretValue, nil
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

/******************************************* SecretFlow Methods *******************************************/

func (c *Client) SecretFlow(secretPath string, secretTitle string, separator string) (string, error) {

	secretPath = strings.TrimSpace(secretPath)
	secretTitle = strings.TrimSpace(secretTitle)
	separator = strings.TrimSpace(separator)

	if secretPath == "" {
		return "", errors.New("Please use a valid Path value")
	}

	if secretTitle == "" {
		return "", errors.New("Please use a valid Title value")
	}

	_, err := c.SignAppin()

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	secret, err := c.SecretGetSecretByPath(secretPath, secretTitle, separator)

	if err != nil {
		error_message := err.Error()
		fmt.Println(error_message)
		c.SignOut()
		return "", errors.New(error_message)
	}

	if strings.ToUpper(secret.SecretType) == "FILE" {
		fileSecretContent, err := c.SecretGetFileSecret(secret.Id)
		if err != nil {
			error_message := err.Error()
			fmt.Println(error_message)
			c.SignOut()
			return "", errors.New(error_message)
		}
		return fileSecretContent, nil
	}

	c.SignOut()
	return secret.Password, nil

}

func (c *Client) SecretGetSecretByPath(secretPath string, secretTitle string, separator string) (entities.Secret, error) {

	path := fmt.Sprintf("secrets-safe/secrets?title=%v&path=%v&separator=%v", secretTitle, secretPath, separator)

	body, err := c.httpRequest(path, "GET", bytes.Buffer{})
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

	var SecretObject entities.Secret
	json.Unmarshal(bodyBytes, &SecretObject)
	return SecretObject, nil
}

func (c *Client) SecretGetFileSecret(secretId string) (string, error) {

	path := fmt.Sprintf("secrets-safe/secrets/%v/file/download", secretId)

	body, err := c.httpRequest(path, "GET", bytes.Buffer{})

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

func (c *Client) SignAppin() (entities.User, error) {
	body, err := c.httpRequest("Auth/SignAppin", "POST", bytes.Buffer{})
	if err != nil {
		return entities.User{}, err
	}
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return entities.User{}, err
	}

	var userObject entities.User
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

func (c *Client) httpRequest(path string, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	url := c.requestPath(path)

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

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%v/%v", c.url, path)
}
