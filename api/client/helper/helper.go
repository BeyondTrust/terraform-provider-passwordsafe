package helper

import (
	"terraform-provider-passwordsafe/api/client/entities"
)

func GetFolderByName(folders []entities.Folder, folderName string) entities.Folder {
	for _, folder := range folders {
		if folder.Name == folderName {
			return folder
		}
	}
	return entities.Folder{}
}

func GetSecretByName(secrets []entities.SecretMetadata, secretName string) entities.SecretMetadata {
	for _, secret := range secrets {
		if secret.Title == secretName {
			return secret
		}
	}
	return entities.SecretMetadata{}
}
