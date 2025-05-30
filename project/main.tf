terraform {
  required_providers {
    passwordsafe = {
      #source  = "providers/beyondtrust/passwordsafe"
      source  = "beyondtrust/passwordsafe"
      version = "1.0.6"
    }
  }
}

// this provider definition combines providerSdkv2 and providerFramework.
provider "passwordsafe" {
  api_key                         = var.api_key
  client_id                       = var.client_id
  client_secret                   = var.client_secret
  url                             = var.url
  api_version                     = var.api_version
  api_account_name                = var.api_account_name
  verify_ca                       = true
  client_certificates_folder_path = var.client_certificates_folder_path
  client_certificate_name         = var.client_certificate_name
  client_certificate_password     = var.client_certificate_password
}


// providerSdkv2

data "passwordsafe_managed_account" "manage_account_01" {
  system_name  = "system01"
  account_name = "managed_account01"
}

output "manage_account_01" {
  value = data.passwordsafe_managed_account.manage_account_01.value
}

