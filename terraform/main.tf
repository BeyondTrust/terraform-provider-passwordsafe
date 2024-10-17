terraform {
  required_providers {
    passwordsafe = {
      source = "providers/beyondtrust/passwordsafe"
      version = "1.0.1"
    }
  }
}

provider "passwordsafe" {
  api_key = "${var.api_key}"
  client_id = "${var.client_id}"
  client_secret = "${var.client_secret}"
  url = "${var.url}"
  api_account_name = "${var.api_account_name}"
  verify_ca = false
  client_certificates_folder_path = "${var.client_certificates_folder_path}"
  client_certificate_name = "${var.client_certificate_name}"
  client_certificate_password = "${var.client_certificate_password}"
}


data "passwordsafe_managed_account" "manage_account_01" {
  system_name = "server01"
  account_name = "managed_account_01"
}

output "manage_account_01" {
  value = "${data.passwordsafe_managed_account.manage_account_01.value}"
}

data "passwordsafe_managed_account" "manage_account_02" {
  system_name = "server01"
  account_name = "managed_account_02"
}

output "manage_account_02" {
  value = "${data.passwordsafe_managed_account.manage_account_02.value}"
}

data "passwordsafe_managed_account" "manage_account_03" {
  system_name = "server01"
  account_name = "managed_account_03"
}

output "manage_account_03" {
  value = "${data.passwordsafe_managed_account.manage_account_03.value}"
}


data "passwordsafe_managed_account" "manage_account_04" {
  system_name = "server01"
  account_name = "managed_account_04"
}

output "manage_account_04" {
  value = "${data.passwordsafe_managed_account.manage_account_04.value}"
}


data "passwordsafe_managed_account" "manage_account_05" {
  system_name = "server01"
  account_name = "managed_account_05"
}

output "manage_account_05" {
  value = "${data.passwordsafe_managed_account.manage_account_05.value}"
}


data "passwordsafe_managed_account" "manage_account_06" {
  system_name = "server01"
  account_name = "managed_account_06"
}

output "manage_account_06" {
  value = "${data.passwordsafe_managed_account.manage_account_06.value}"
}


data "passwordsafe_managed_account" "manage_account_07" {
  system_name = "server01"
  account_name = "managed_account_07"
}

output "manage_account_07" {
  value = "${data.passwordsafe_managed_account.manage_account_07.value}"
}

data "passwordsafe_managed_account" "manage_account_08" {
  system_name = "server01"
  account_name = "managed_account_08"
}

output "manage_account_08" {
  value = "${data.passwordsafe_managed_account.manage_account_08.value}"
}


data "passwordsafe_managed_account" "manage_account_09" {
  system_name = "server01"
  account_name = "managed_account_09"
}

output "manage_account_09" {
  value = "${data.passwordsafe_managed_account.manage_account_09.value}"
}

data "passwordsafe_managed_account" "managed_account_10" {
  system_name = "server01"
  account_name = "managed_account_10"
}

output "managed_account_10" {
  value = "${data.passwordsafe_managed_account.managed_account_10.value}"
}


data "passwordsafe_secret" "secret_text" {
  path = "local/folder"
  title = "my_credential"
}

output "secret_text" {
  value = "${data.passwordsafe_secret.secret_text.value}"
}