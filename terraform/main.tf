terraform {
  required_providers {
    passwordsafe = {
      source = "providers/beyondtrust/passwordsafe"
      version = "1.0.0"
    }
  }
}

provider "passwordsafe" {
  apikey = "${var.api_key}"
  url    = "${var.url}"
  accountname = "${var.account_name}"
}

data "passwordsafe_managed_account" "manage_account" {
  system_name = "Computer01"
  account_name = "User05"
}


output "manage_account" {
  value = "${data.passwordsafe_managed_account.manage_account.value}"
}


