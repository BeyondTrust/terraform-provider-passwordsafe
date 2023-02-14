// Copyright 2023 BeyondTrust. All rights reserved.
// Package entities implements DTO's used by Beyondtrust Secret Safe API.
package entities

type User struct {
	UserId       int    `json:"UserId"`
	EmailAddress string `json:"EmailAddress"`
	UserName     string `json:"UserName"`
	Name         string `json:"Name"`
}

type ManagedAccount struct {
	SystemId  int
	AccountId int
}

type Secret struct {
	Id         string
	Title      string
	Password   string
	SecretType string
}
