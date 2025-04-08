terraform {
  required_providers {
    passwordsafe = {
      source  = "providers/beyondtrust/passwordsafe"
      version = "1.0.1"
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
  system_name  = "server01"
  account_name = "managed_account_01"
}

output "manage_account_01" {
  value = data.passwordsafe_managed_account.manage_account_01.value
}

data "passwordsafe_managed_account" "manage_account_02" {
  system_name  = "server01"
  account_name = "managed_account_02"
}

output "manage_account_02" {
  value = data.passwordsafe_managed_account.manage_account_02.value
}

data "passwordsafe_managed_account" "manage_account_03" {
  system_name  = "server01"
  account_name = "managed_account_03"
}

output "manage_account_03" {
  value = data.passwordsafe_managed_account.manage_account_03.value
}


data "passwordsafe_managed_account" "manage_account_04" {
  system_name  = "server01"
  account_name = "managed_account_04"
}

output "manage_account_04" {
  value = data.passwordsafe_managed_account.manage_account_04.value
}


data "passwordsafe_managed_account" "manage_account_05" {
  system_name  = "server01"
  account_name = "managed_account_05"
}

output "manage_account_05" {
  value = data.passwordsafe_managed_account.manage_account_05.value
}


data "passwordsafe_managed_account" "manage_account_06" {
  system_name  = "server01"
  account_name = "managed_account_06"
}

output "manage_account_06" {
  value = data.passwordsafe_managed_account.manage_account_06.value
}


data "passwordsafe_managed_account" "manage_account_07" {
  system_name  = "server01"
  account_name = "managed_account_07"
}

output "manage_account_07" {
  value = data.passwordsafe_managed_account.manage_account_07.value
}

data "passwordsafe_managed_account" "manage_account_08" {
  system_name  = "server01"
  account_name = "managed_account_08"
}

output "manage_account_08" {
  value = data.passwordsafe_managed_account.manage_account_08.value
}


data "passwordsafe_managed_account" "manage_account_09" {
  system_name  = "server01"
  account_name = "managed_account_09"
}

output "manage_account_09" {
  value = data.passwordsafe_managed_account.manage_account_09.value
}

data "passwordsafe_managed_account" "managed_account_10" {
  system_name  = "server01"
  account_name = "managed_account_10"
}

output "managed_account_10" {
  value = data.passwordsafe_managed_account.managed_account_10.value
}


data "passwordsafe_secret" "secret_text" {
  path  = "local/folder"
  title = "my_credential"
}

output "secret_text" {
  value = data.passwordsafe_secret.secret_text.value
}


resource "passwordsafe_managed_account" "my_managed_account" {
  system_name  = "system_integration_test"
  account_name = "managed_account_Test"
  password     = "MyTest101*!"
  api_enabled  = true
}


resource "passwordsafe_credential_secret" "my_credenial_secret" {
  folder_name = "folder1"
  title       = "Credential_Secret_from_Terraform"
  description = "my credential secret description"
  username    = "my_user_name"
  password    = "password_content"
  owner_type  = "User"
  notes       = "My Notes"
  group_id    = 1
}


resource "passwordsafe_text_secret" "my_text_secret" {
  folder_name = "folder1"
  title       = "Text_Secret_from_Terraform"
  description = "my text secret description"
  owner_type  = "User"
  text        = "password_text"
  notes       = "My notes"
  group_id    = 1
}


resource "passwordsafe_file_secret" "my_file_secret" {
  folder_name  = "folder1"
  title        = "File_Secret_from_Terraform"
  description  = "my file secret description"
  owner_type   = "User"
  file_content = file("test_secret.txt")
  file_name    = "my_secret.txt"
  notes        = "My notes"
  group_id     = 1
}

resource "passwordsafe_folder" "my_folder" {
  parent_folder_name = "folder1"
  name               = "my_new_folder_mame"
}

resource "passwordsafe_safe" "my_safe" {
  name        = "my_new_safe_mame"
  description = "my_safe_description"
}


// providerFramework (passwordsafe_managed_acccount_ephemeral, passwordsafe_secret_ephemeral)

ephemeral "passwordsafe_managed_acccount_ephemeral" "managed_account" {
  system_name  = "system01"
  account_name = "managed_account01"
}


ephemeral "passwordsafe_secret_ephemeral" "secret" {
  path  = "oauthgrp"
  title = "ephemeral_secret_title1"
}

resource "passwordsafe_workgroup" "workgroup" {
  name = "workgroup name"
}

resource "passwordsafe_asset_by_workgroup_name" "asset_by_workgroup_name" {
  work_group_name  = "Wrokgroup Name"
  ip_address       = "10.20.30.40"
  asset_name       = "Prod_Server_03"
  dns_name         = "server01.company.com"
  domain_name      = "company.com"
  mac_address      = "00:1A:2B:3C:4D:5E"
  asset_type       = "Windows Server"
  description      = "Production Windows Server hosting critical applications"
  operating_system = "Windows Server 2022"
}

resource "passwordsafe_asset_by_workgroup_id" "asset_by_workgroup_id" {
  work_group_id    = "28"
  ip_address       = "10.20.30.40"
  asset_name       = "Prod_Server_03"
  dns_name         = "server01.company.com"
  domain_name      = "company.com"
  mac_address      = "00:1A:2B:3C:4D:5E"
  asset_type       = "Windows Server"
  description      = "Production Windows Server hosting critical applications"
  operating_system = "Windows Server 2022"
}


resource "passwordsafe_database" "database" {
  asset_id            = "1"
  platform_id         = 10
  instance_name       = "primary-db-instance"
  is_default_instance = false
  port                = 5432
  version             = "13.3"
  template            = "standard-template"
}

resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset Description"
  port                                   = 5432
  timeout                                = 30
  ssh_key_enforcement_mode               = 1
  password_rule_id                       = 0
  dss_key_rule_id                        = 0
  login_account_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 90
  auto_management_flag                   = false
  functional_account_id                  = 0
  elevation_command                      = "sudo su"
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "last"
  change_frequency_days                  = 30
  change_time                            = "02:00"
  remote_client_type                     = "EPM"
  application_host_id                    = 0
  is_application_host                    = false
}


resource "passwordsafe_managed_system_by_workgroup" "managed_system_by_workgroup" {
  workgroup_id                           = "55"
  entity_type_id                         = 1
  host_name                              = "example-host"
  ip_address                             = "222.222.222.22"
  dns_name                               = "example.local"
  instance_name                          = "example-instance"
  is_default_instance                    = true
  template                               = "example-template"
  forest_name                            = "example-forest"
  use_ssl                                = false
  platform_id                            = 2
  netbios_name                           = "EXAMPLE"
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_workgroup Description"
  port                                   = 5432
  timeout                                = 30
  ssh_key_enforcement_mode               = 0
  password_rule_id                       = 0
  dss_key_rule_id                        = 0
  login_account_id                       = 0
  account_name_format                    = 1
  oracle_internet_directory_id           = "example-dir-id"
  oracle_internet_directory_service_name = "example-service"
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 30
  auto_management_flag                   = false
  functional_account_id                  = 0
  elevation_command                      = "sudo su -"
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "last"
  change_frequency_days                  = 7
  change_time                            = "02:00"
  access_url                             = "https://example.com"
  remote_client_type                     = "ssh"
  application_host_id                    = 5001
  is_application_host                    = false
}


resource "passwordsafe_managed_system_by_database" "managed_system_by_database" {
  database_id                            = "2"
  contact_email                          = "admin@example.com"
  description                            = "Managed system for example DB"
  timeout                                = 30
  password_rule_id                       = 101
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 45
  auto_management_flag                   = true
  functional_account_id                  = 1234
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "xdays"
  change_frequency_days                  = 15
  change_time                            = "03:00"
}
