package main

import (
	"fmt"
	"strconv"
	"terraform-provider-passwordsafe/api/client"
	"terraform-provider-passwordsafe/config"
)

func main_test() {
	apiKey := config.PS_API_KEY
	accountName := config.PS_ACCOUNT_NAME
	url := config.PS_URL
	apiClient := client.NewClient(url, apiKey, accountName)

	secret, err := apiClient.ManageAccountFlow("Computer01", "User04")

	if err != nil {
		fmt.Printf("Error in ManageAccountFlow: %v", err.Error())
	}

	secret, _ = strconv.Unquote(secret)
	fmt.Println(secret)

}
