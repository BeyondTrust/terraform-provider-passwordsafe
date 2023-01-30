package entities

type User struct {
	UserId       int    `json:"UserId"`
	EmailAddress string `json:"EmailAddress"`
	UserName     string `json:"UserName"`
	Name         string `json:"Name"`
}

type Folder struct {
	Id          string
	Name        string
	ParentId    string
	Description string
	UserGroupId int
}

type SecretMetadata struct {
	Id         string
	Title      string
	FolderPath string
	SecretType string
}

type TextSecret struct {
	Id    string
	Title string
	Text  string
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
