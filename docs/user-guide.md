The Password Safe Terraform provider is a Terraform integration that enables the use of Password Safe secrets management capabilities with Terraform.

Terraform configuration files can be configured to retrieve secrets from Password Safe.

Permissions to access secrets in Password Safe can be granted to specific accounts within BeyondInsight.

> ℹ️ 
> 
> The Password Safe Terraform Provider is available in the [Terraform Registry](https://registry.terraform.io/providers/BeyondTrust/passwordsafe/latest).

## Overview

The Password Safe Terraform provider is named terraform-provider-passwordsafe. It is available for the following Operating systems: Windows, Linux and MacOS. 

The provider enables using Password Safe secrets management capabilities with Terraform. The provider is invoked by using HashiCorp Configuration Language (HCL) which tells the provider to reach out (using APIs) to the Password Safe service and utilize its secrets management capabilities.

1. The end user writes terraform configuration that sets PS variables for the PS provider and retrieves the API secret for authentication to a 3rd party provider. All configuration is run on Terraform plan/apply.
2. Terraform plan/apply invokes the Password Safe provider to retrieve the secret.
3. Terraform plan/apply also invokes the third-party provider, using the Password Safe secret for authentication to the third-party service.

> ℹ️ 
> 
> This integration requires Password Safeversion 23.1 or higher for Secrets Safe secrets.
> 
> This integration uses Transport Layer Security (TLS) 1.2 protocol.

> ℹ️ 
> 
> For more information, see [The Core Terraform Workflow](https://developer.hashicorp.com/terraform/intro/core-workflow).

![A workflow diagram showing how Terraform interacts with Password Safe and Secrets Safe. Terraform Core communicates with plugins over RPC. The plugins include the PS Terraform Provider and a third‑party Terraform provider. The PS Terraform Provider connects over HTTPS to Password Safe to retrieve managed accounts and to Secrets Safe for credential files. A separate third‑party service connects to the third‑party provider. The diagram labels steps for configuration, provider calls, and secret retrieval.](https://files.readme.io/9545534e3f6d6199e2dc600286b60cc1fbaf42db75283af017a6f95b874c7892-image.png)

## Prerequisites

- Requires Password Safe 24.3 or later for Secrets Safe
- Uses TLS 1.2 protocol
- Terraform v1.10 is required.

## Configure Password Safe to allow Terraform files to retrieve secrets

The following sections outline how to set up the required authorization in Password Safe to allow Terraform HCL files to retrieve secrets, and then shows examples of Terraform HCL files for retrieving secrets.

> ℹ️ 
> 
> The following instructions are meant to be a quick start guide to help with Password Safe setup.

### Password Safe Authentication Methods

The provider must be configured with one of next authentication methods:

1. API key and a user account to access target secrets in Password Safe. 
2. Client_id/client_secret using OAuth method to access target secrets in Password Safe. 

> ℹ️ 
> 
> For Authentication: Use api\_key or client\_id/client_secret, when api\_key is empty so client\_id and client\_secret are used to authenticate using API OAuth method.

This is accomplished by the following steps:

#### Set up API key and user account

1. Create an API registration with API key or OAuth method in BeyondInsight.
2. Create or use an existing Secrets Safe group.
3. Create or use an existing BeyondInsight user.
4. Add the API registration to the group.
5. Add the user to the group.
6. Add the Secrets Safe feature to the group.

#### Managed accounts setup

1. Create or use an existing access policy that has the **View Password Auto Approve** option set.
2. Add the **All Managed Accounts Smart Group** to the BeyondInsight group.
3. Add the access policy to the **All Managed Accounts Smart Group** role, and ensure that both requestor and approver are set.
4. Create or use an existing managed system.
5. Create or use an existing managed account associated with the managed system.
6. Configure the managed account with the **API Enabled** and **Max Concurrent Requests Unlimited** options selected.

## Configure Terraform

> 🚧 Important information
> 
> Terraform state files and plan files contain sensitive information. Ensure best practices are followed for securing files.

### Plugin directory setup

The provider must be placed in the terraform plugin directory. The plugin directory can be customized by the end user.

The default unpacked layout, implied local mirror directory is:

- **Windows**: %APPDATA%/terraform.d/plugins/providers/beyondtrust/passwordsafe/1.0.0/windows_amd64
- **Mac OS X**: $HOME/.terraform.d/plugins/providers/beyondtrust/passwordsafe/1.0.0/darwin_amd64
- **Linux systems**: $HOME/.terraform.d/plugins/providers/beyondtrust/passwordsafe/1.0.0/linux_amd64

> ⚠️ Warning
> 
> Remember to set the executable bit on the provider chmod +x

> ℹ️ 
> 
> For more information, see [Implied Local Mirror Directories](https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories).

The path must end with the provider (terraform-provider-passwordsafe_1.0.0) VERSION, and a string like 1.0.0. The TARGET specifies a particular target platform using a format like darwin_amd64, linux_amd64, windows_amd64.

### Configure provider

Configure the provider in your HCL files.

> ✅ **Example**
> 
> Provider Configuration Example:

```
terraform {
  required_providers {
    passwordsafe = {
      source = "providers/beyondtrust/passwordsafe"
      version = "1.0.0"
    }
  }
}
# configure the Password Safe provider
provider "passwordsafe" {
  # Note about Authentication: (Use api_key or client_id/client_secret, 
  # when api_key is empty so client_id and client_secret will be used 
  # to authenticate using API OAuth method).
  api_key = "${var.api_key}"
  client_id = "${var.client_id}"
  client_secret = "${var.client_secret}"
  url = "${var.url}"
  api_version = "${var.api_version}"
  api_account_name = "${var.api_account_name}"
  verify_ca = true
  client_certificates_folder_path = "${var.client_certificates_folder_path}"
  client_certificate_name = "${var.client_certificate_name}"
  client_certificate_password = "${var.client_certificate_password}"
}
```

#### Provider Arguments

- **url**: The URL for the Password Safe instance used to request a secret..
- **api_version**: The recommended version is 3.1. If no version is specified, the default API version 3.0 will be use.
- **api_account_name**: The user name for the api request to the Password Safe instance. For use when authenticating with an api key.
- **api_key**: The api key for making requests to the Password Safe instance. For use when authenticating to Password Safe. 
- **client_id**: API OAuth Client ID. 
- **client_secret**: API OAuth Client Secret.
- **client_certificates_folder_path**: (optional) The path to the Client Certificate associated with the Password Safe instance for use when authenticating with an API key using a Client Certificate.
- **client_certificate_password**: (optional) The password associated with the Client Certificate. For use when authenticating with an API key using a Client Certificate.
- **client_certificate_name**: (optional) The name of the Client Certificate for use when authenticating with an API key using a Client Certificate.
- **verify_ca**: (optional) Indicates whether to verify the certificate authority on the Password Safe instance. For use when authenticating to Password Safe.

### BeyondTrust Terraform available datasources

Data is retrieved from the Password Safe API. The response format depends on the API user configured by the Terraform provider.

List of available datasources:

**Datasources to retrieve a single record**

- **passwordsafe_managed_account**: retrieve a managed account secret.
- **passwordsafe_secret**: retrieve a secret
- **passwordsafe_managed_acccount_ephemeral**: retrieve a managed account secret. (ephemeral datasource)
- **passwordsafe_secret_ephemeral**: retrieve a secret. (ephemeral datasource)

**Datasources to retrieve lists**

- **passwordsafe_asset_datasource**: get list of assets.
- **passwordsafe_functional_account_datasource**: get list of functional accounts.
- **passwordsafe_database_datasource**: get list of databases.
- **passwordsafe_folder_datasource:** get list of folders.
- **passwordsafe_safe_datasource**: get list of safes datasources.
- **passwordsafe_workgroup_datasource**: get list of workgroups
- **passwordsafe_platform_datasource**: get list of platforms.
- **passwordsafe_managed_account_datasource**: get list of managed accounts.
- **passwordsafe_managed_system_datasource**: get list of managed systems.

```
// assets_list
data "passwordsafe_asset_datasource" "assets_list" {
  parameter = "1" // could be workgroup id / workgroup name.
}
output "assets_list" {
  value = data.passwordsafe_asset_datasource.assets_list.assets[0].asset_name
}
// functional_accounts_list
data "passwordsafe_functional_account_datasource" "functional_accounts_list" {
}
output "functional_account_list" {
  value = [
    for acc in data.passwordsafe_functional_account_datasource.functional_accounts_list.accounts : acc.account_name
    if acc.account_name == "svc-monitoring"
  ][0]
}
locals {
  functional_account_name = [
    for acc in data.passwordsafe_functional_account_datasource.functional_accounts_list.accounts : acc.account_name
    if acc.account_name == "svc-monitoring"
  ][0]
}
// databases_list
data "passwordsafe_database_datasource" "databases_list" {
}
output "database_list" {
  value = data.passwordsafe_database_datasource.databases_list.databases[1].instance_name
}
// folders_list
data "passwordsafe_folder_datasource" "folders_list" {
}
output "folders_list" {
  value = data.passwordsafe_folder_datasource.folders_list.folders[0].name
}
// safes_list
data "passwordsafe_safe_datasource" "safes_list" {
}
output "safes_list" {
  value = data.passwordsafe_safe_datasource.safes_list.safes[0].name
}
// workgroups_list
data "passwordsafe_workgroup_datasource" "workgroups_list" {
}
output "workgroups_list" {
  value = data.passwordsafe_workgroup_datasource.workgroups_list.workgroups[0].name
}
// plarforms_list
data "passwordsafe_platform_datasource" "plarforms_list" {
}
output "plarforms_list" {
  value = data.passwordsafe_platform_datasource.plarforms_list.platforms[0].name
}
// managed_accounts_list
data "passwordsafe_managed_account_datasource" "managed_accounts_list" {
}
output "managed_accounts_list" {
  value = data.passwordsafe_managed_account_datasource.managed_accounts_list.managed_accounts[0].account_name
}
// managed_system_list
data "passwordsafe_managed_system_datasource" "managed_system_list" {
}
output "managed_system_list" {
  value = data.passwordsafe_managed_system_datasource.managed_system_list.managed_systems[0].system_name
}
```

### Beyondtrust terraform available resources

Sends data to Password Safe API.

List of Available Resources:

- **passwordsafe_managed_account**: create managed account.
- **passwordsafe_credential_secret**: create credential secret.
- **passwordsafe_text_secret**: create text secret.
- **passwordsafe_file_secret**: create file secret.
- **passwordsafe_folder**: create folder.
- **passwordsafe_safe**: create safe.
- **passwordsafe_workgroup**: create workgroup.
- **passwordsafe_asset_by_workgroup_name**: create asset using workgroup name.
- **passwordsafe_asset_by_workgroup_id**: create asset using workgroup id.
- **passwordsafe_database**: create database
- **passwordsafe_managed_system_by_asset**: create managed system using asset id.
- **passwordsafe_managed_system_by_workgroup**: create managed system using workgroup id. 
- **passwordsafe_managed_system_by_database**: create database using database id.
- **passwordsafe_functional_account**: create functional account. (work in progress)

### Configure secrets retrieval

Configure data sources to retrieve secrets from Password Safe. The first data block in the example below is _passwordsafe\_managed\_account_, used to retrieve managed account secrets. The second type of data source is _passwordsafe\_secret_, used to retrieve credentials, text, and file secrets.

The **decrypt**  parameter is optional , When set to true, the decrypted password field is returned. When set to **false**, the password field is omitted. This option applies only to secret retrieval type. The default value is set to true if not specified.

> ✅ Example

```
# retrieve a managed account secret
data "passwordsafe_managed_account" "manage_account" {
  system_name = "ServerStandard"
  account_name = "serveruser1"
}
output "manage_account" {
  value = "${data.passwordsafe_managed_account.manage_account.value}"
}
# retrieve a secrets safe credential
data "passwordsafe_secret" "secret_credential" {
  path = "folder1*folder2*folder3/folder4"
  title = "credLevel6"
  # separator is optional, and its default value is '/'.
  separator = "*"
  # decrypt is optional, and its default value is true.
  decrypt = "true"
}
output "secret_credential" {
  value = "${data.passwordsafe_secret.secret_credential.value}"
}
# retrieve a secrets safe file
data "passwordsafe_secret" "secret_file" {
  path = "folder1"
  title = "RootFile"
}
 output "secret_file" {
  value = "${data.passwordsafe_secret.secret_file.value}"
}
```

### Configure Secrets Retrieval using ephemeral resources

The first data block in the example below, _passwordsafe_managed_account_ephemeral_, is used to retrieve managed account secrets using an ephemeral resource.

The second data source, _passwordsafe_secret_ephemeral_, is used to retrieve credentials, text, and file secrets using an ephemeral resource. These resources do not store any data in state files.

> ℹ️ 
> 
> For more information, see [Terraform Ephemeral Resources.](https://developer.hashicorp.com/terraform/language/resources/ephemeral)

```
ephemeral "passwordsafe_managed_acccount_ephemeral" "managed_account" {
  system_name = "system01"
  account_name = "managed_account01"
}
ephemeral "passwordsafe_secret_ephemeral" "secret" {
  path = "oauthgrp"
  title = "ephemeral_secret_title"

```

> ℹ️ 
> 
> Ephemeral resources are available in Terraform v1.10 and later.

### Create Secrets, folder and safes

 You can create three types of secrets: 

- Credential
- Text
- File

You need to define the parent folder name. If the parent folder does not exist, you must create it before creating the secret.

#### Create Credential Secret

> ✅ **Example**

```
# create credential secret
# Supported versions: 3.0, 3.1, and 3.2.
resource "passwordsafe_credential_secret" "my_credenial_secret" {
  folder_name = "folder1"
  title = "Credential_Secret_from_Terraform"
  username = "my_user_name"
  password = "password_content"
  owner_type = "User"
  
  # optional attributes from here on
  description = "my credential secret description"
  owner_id = 1
  notes = "My Notes"
  password_rule_id = 1
  
  # the API user will be added as the owner by default
  # version 3.0 (using owner_id)
  owners {
    owner_id = 1
    owner = "owner"
    email = "user@beyondtrust.com"
  }
  
  # the API user will be added as the owner by default
  # version 3.1, 3.2 (using user_id)
  owners {
    user_id = 1
    name = "owner"
    email = "user@beyondtrust.com"
  }

  urls {
    id = 1
    credential_id = 1
    url = "https://www.beyondtrust.com/"
  }
}
```

**Parameters Description**

- The folder name indicates where the secret is created.
- Max string length for description and password is 256.
- Max string length for notes is 4000.
- Max string length for URL is 2048.
- Required: Title, username, password.
- A password or a PasswordRuleID is required.
  - If a PasswordRuleID is passed in, then a password is generated (based on the Password Policy defined by the PasswordPolicyID).
  - If a password is passed in instead, the same behavior is followed (using that as the password).

#### Create Text Secret

> ✅ **Example**

```
# create text secret
# Supported versions: 3.0, 3.1, and 3.2.
resource "passwordsafe_text_secret" "my_text_secret" {
  folder_name = "folder1"
  title = "Text_Secret_from_Terraform"
  text = "password_text"
  owner_type = "User"
  # optional attributes from here on
  description = "my text secret description"
  owner_id = 1
  notes = "My notes"
  # the API user will be added as the owner by default
  # version 3.0 (using owner_id)
  owners {
    owner_id = 1
    owner = "owner"
    email = "user@beyondtrust.com"
  }
  # the API user will be added as the owner by default
  # version 3.1, 3.2 (using user_id)
  owners {
    user_id = 1
    name = "owner"
    email = "user@beyondtrust.com"
  }
  urls {
    id = 1
    credential_id = 1
    url = "https://www.beyondtrust.com/"
  }
}
```

**Parameters Description**

- The folder name indicates where the secret is created.
- Max string length for Title and Description is 256.
- Max string length for text is 4096.
- Max string length for notes is 4000.
- Max string length for URL is 2048.
- Required: Title, FolderId.

#### Create File Secret

> ✅ **Example**

```
# create file secret
# Supported versions: 3.0, 3.1, and 3.2.
resource "passwordsafe_file_secret" "my_file_secret" {
  folder_name = "folder1"
  title = "File_Secret_from_Terraform"
  file_content = file("test_secret.txt")
  file_name = "my_secret.txt"
  owner_type = "User"
  
  # optional attributes from here on
  description = "my file secret description"
  owner_id = 1
  notes= "My notes"
  
  # the API user will be added as the owner by default
  # version 3.0 (using owner_id)
  owners {
    owner_id = 1
    owner = "owner"
    email = "user@beyondtrust.com"
  }

  # the API user will be added as the owner by default
  # version 3.1, 3.2 (using user_id)
  owners {
    user_id = 1
    name = "owner"
    email = "user@beyondtrust.com"
  }
  
  urls {
    id = 1
    credential_id = 1
    url = "https://www.beyondtrust.com/"
  }
}
```

**Parameters Description**

- The folder name indicates where the secret is created.
- Max string length for Title, Description, and FileName is 256.
- Max string length for notes is 4000.
- Max string length for URL is 2048.
- Max file size is 5 <Glossary>MB</Glossary>. Size must be greater than 0 MB.
- Required: Title, FolderId, Filename.

#### Create Folder

> ✅ **Example**

```
# create folder
resource "passwordsafe_folder" "my_folder" {
  parent_folder_name = "folder1"
  name = "my_new_folder_mame"
  # optional attributes from here on
  description = "New Folder Description"
  user_group_id = 1
}
```

**Parameters Description**

- **parent_folder_name**: (required) Parent folder where new folder will be created.
- **name**: (required) New folder name.
- **description**: A description of the folder. 
- **user_group_id**: User group id.

#### Create Safe

> ✅ **Example**

```
# create safe
resource "passwordsafe_safe" "my_safe" {
  name = "my_new_safe_mame"
  description="my_safe_description" #optional
}
```

**Parameters Description**

- **name**: (required) New safe name.
- **description**: A description of the safe.

#### Create Managed Account

> ✅ **Example**

```
# create managed account
resource "passwordsafe_managed_account" "my_managed_account" {
  system_name = "system_integration_test"
  account_name = "managed_account_Test"
  password = "MyTest101*!"
  # optional attributes from here on
  domain_name                       = "example.com"
  user_principal_name               = "user123@example.com"
  sam_account_name                  = "USER123"
  distinguished_name                = "CN=User123,OU=Users,DC=example,DC=com"
  private_key                       = "/path/to/private/key.pem"
  passphrase                        = "mysecretpassphrase"
  password_fallback_flag            = true
  login_account_flag                = false
  description                       = "This is a test user account"
  password_rule_id                  = 1
  api_enabled                       = true
  release_notification_email        = "notify@example.com"
  change_services_flag              = true
  restart_services_flag             = false
  change_tasks_flag                 = true
  release_duration                  = 30
  max_release_duration              = 90
  isa_release_duration              = 60
  max_concurrent_requests           = 10
  auto_management_flag              = true
  dss_auto_management_flag          = false
  check_password_flag               = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag   = true
  change_frequency_type             = "first"
  change_frequency_days             = 30
  change_time                       = "03:00"
  next_change_date                  = "2025-01-01"
  use_own_credentials               = false
  workgroup_id                      = 2
  change_windows_auto_logon_flag    = true
  change_com_plus_flag              = false
  change_dcom_flag                  = true
  change_scom_flag                  = false
  object_id                         = "12345-abcde-67890"
}
```

**Parameters Description**

- **system_name**  (required): The name of the system, Max string length is 245.
- **account_name**  (required): The name of the account. Must be unique on the system. Max string length is 245.
- **password**  (required if auto_management_flag is false): The account password.
- **domain_name**  (optional): This can be given but it must be exactly the same as the directory. If empty or null, it is automatically populated from the parent managed system/directory. Max string length is 50.
- **user_principal_name**  (required for Active Directory and Entra ID managed systems only): The Active Directory user principal name. Max string length is 500.
- **sam_account_name**  (required for Active Directory managed systems, optional for Entra ID managed systems): The Active Directory SAM account name (Maximum 20 characters). Max string length is 20.
- **distinguished_name**  (required for LDAP Directory managed systems only): The LDAP distinguished name. Max string length is 1000.
- **private_key**: DSS private key. Can be set if platform.dss_flag is true.
- **passphrase**  (required when private_key is an encrypted DSS key): DSS passphrase. Can be set if platform.dss_flag is true.
- **password_fallback_flag** (default: false): True if failed DSS authentication can fall back to password authentication, otherwise false. Can be set if platform.dss_flag is true.
- **login_account_flag**: True if the account should use the managed system login account for SSH sessions, otherwise false. Can be set when the managed_system.login_account_id is set.
- **description**: A description of the account. Max string length is 1024.
- **password_rule_id** (default: 0): ID of the password rule assigned to this managed account.
- **api_enabled** (default: false): True if the account can be requested through the API, otherwise false.
- **release_notification_email**: Email address used for notification emails related to this managed account. Max string length is 255.
- **change_services_flag** (default: false): True if services run as this user should be updated with the new password after a password change, otherwise false.
- **restart_services_flag** (default: false): True if services should be restarted after the run as password is changed (change_services_flag), otherwise false.
- **change_tasks_flag** (default: false): True if scheduled tasks run as this user should be updated with the new password after a password change, otherwise false.
- **release_duration** (minutes: 1-525600, default: 120): Default release duration.
- **max_release_duration** (minutes: 1-525600, default: 525600): Default maximum release duration.
- **isa_release_duration** (minutes: 1-525600, default: 120): Default Information Systems Administrator (ISA) release duration.
- **max_concurrent_requests** (0-999, 0 is unlimited, default: 1): Maximum number of concurrent password requests for this account.
- **auto_management_flag** (default: false): True if password auto-management is enabled, otherwise false.
- **dss_auto_management_flag** (default: false): True if DSS key auto-management is enabled, otherwise false. If set to true, and no private_key is provided, immediately attempts to generate and set a new public key on the server. Can be set if platform.dss_auto_management_flag is true.
- **check_password_flag** (default: false): True to enable password testing, otherwise false.
- **change_password_after_any_release_flag** (default: false): True to change passwords on release of a request, otherwise false.
- **reset_password_on_mismatch_flag** (default: false): True to queue a password change when scheduled password test fails, otherwise false.
- **change_frequency_type**  (default: first): The change frequency for scheduled password changes:
  - **first**: Changes scheduled for the first day of the month.
  - **last**: Changes scheduled for the last day of the month.
  - **xdays**: Changes scheduled every x days (change_frequency_days).
- **change_frequency_days** (days: 1-999): When change_frequency_type is xdays, password changes take place this configured number of days.
- **change_time** (24hr format: 00:00-23:59, default: 23:30): UTC time of day scheduled password changes take place.
- **next_change_date** (date format: YYYY-MM-DD): UTC date when next scheduled password change occurs. If the next_change_date + change_time is in the past, password change occurs at the nearest future change_time.
- **use_own_credentials** (version 3.1+): True if the current account credentials should be used during change operations, otherwise false.
- **change_iis_app_pool_flag** (version 3.2 only): True if IIS application pools run as this user should be updated with the new password after a password change, otherwise false.
- **restart_iis_app_pool_flag** (version 3.2 only): True if IIS application pools should be restarted after the run as password is changed (change_iis_app_pool_flag), otherwise false.
- **workgroup_id**: ID of the assigned Workgroup.
- **change_windows_auto_logon_flag** (default: false): True if Windows Auto Logon should be updated with the new password after a password change, otherwise false.
- **change_com_plus_flag** (default: false): True if COM+ Apps should be updated with the new password after a password change, otherwise false.
- **change_dcom_flag** (default: false): True if DCOM Apps should be updated with the new password after a password change, otherwise false.
- **change_scom_flag** (default: false): True if SCOM Identities should be updated with the new password after a password change, otherwise false.
- **object_id** (required when platform.requires_object_id is true): ObjectID of the account (if applicable). Max string length is 36.

#### Create Workgroup

```
# create workgroup
resource "passwordsafe_workgroup" "workgroup" {
  name = "workgroup_name"
}
```

**Parameters Description**

- **Organization ID**: (optional) The ID of the organization in which to place the new Workgroup. If empty, the Workgroup is placed in the default organization.
- **Name**: The name of the Workgroup. Max string length is 256.

#### Create Asset using workgroup name

```
# create asset by workgroup name
resource "passwordsafe_asset_by_workgroup_name" "asset_by_workgroup_name" {
  work_group_name  = passwordsafe_workgroup.workgroup.name
  ip_address       = "10.20.30.40"
  asset_name       = "Prod_Server_03"
  dns_name         = "server01.company.com"
  domain_name      = "company.com"
  mac_address      = "00:1A:2B:3C:4D:5E"
  asset_type       = "Windows Server"
  description      = "Production Windows Server hosting critical applications"
  operating_system = "Windows Server 2022"
}
```

**Parameters Description**

- **Workgroup name**: (required) Workgroup name.
- **IPAddress**: (required) Asset IP address. Max string length is 45.
- **AssetName**: (optional) Asset name. If not given, a padded IP address is used. Max string length is 128.
- **DnsName**: (optional) Asset DNS name. Max string length is 255.
- **DomainName**: (optional) Asset domain name. Max string length is 64.
- **MacAddress**: (optional) Asset MAC address. Max string length is 128.
- **AssetType**: (optional) Asset type. Max string length is 64.
- **Description**: (optional) Asset description. Only updated if version in the URL is 3.1 or greater. Max string length is 255.
- **OperatingSystem**: (optional) Asset operating system. Max string length is 255.

#### Create Asset using workgroup id

```
# create asset by workgroup id
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
```

**Parameters Description**

- **Workgroup Id**: (required) Workgroup Id.
- **IPAddress**: (required) Asset IP address. Max string length is 45.
- **AssetName**: (optional) Asset name. If not given, a padded IP address is used. Max string length is 128.
- **DnsName**: (optional) Asset DNS name. Max string length is 255.
- **DomainName**: (optional) Asset domain name. Max string length is 64.
- MacAd**d**ress: (optional) Asset MAC address. Max string length is 128.
- **AssetType**: (optional) Asset type. Max string length is 64.
- **Description**: (optional) Asset description. Only updated if version in the URL is 3.1 or greater. Max string length is 255.
- **OperatingSystem**: (optional) Asset operating system. Max string length is 255.

#### Create Database using asset id

```
# create database by asset id
resource "passwordsafe_database" "database" {
  asset_id            = "1"
  platform_id         = 10
  instance_name       = "primary-db-instance"
  is_default_instance = false
  port                = 5432
  version             = "13.3"
  template            = "standard-template"
}
```

**Parameters Description**

- **Asset Id**: (required) Asset Id.
- **PlatformID**: (required) ID of the platform.
- **InstanceName**: Name of the database instance. Required when **IsDefaultInstance** is false. Max string length is 100.
- **IsDefaultInstance**: True if the database instance is the default instance, otherwise false.
- **Port**: (required) The database port.
- **Version**: The database version. Max string value is 20.
- **Template**: The database connection template.

#### Create Managed System using asset id

```
#  create managed system by asset id / API Version: 3.0
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
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
}
# create managed system by asset id / API Version: 3.1
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
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
}
# create managed system by asset id / API Version: 3.2
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
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
```

**Parameters Description**

See available API versions for this resource: [Password Safe APIs ](https://docs.beyondtrust.com/bips/docs/api-reference)

- **Asset ID**: (required) Asset Id.
- **PlatformID**:(required) ID of the managed system platform.
- **ContactEmail**: Max string length is 1000.
- **Description**: Max string length is 255.
- **Port**: (optional) The port used to connect to the host. If null and the related **Platform.PortFlag** is true, Password Safe uses **Platform.DefaultPort** for communication.
- **Timeout**: (seconds, default: 30) Connection timeout. Length of time in seconds before a slow or unresponsive connection to the system fails.
- **SshKeyEnforcementMode**: (default: 0/None) Enforcement mode for SSH host keys.
  - **0**: None.
  - **1**: Auto. Auto accept initial key.
  - **2**: Strict. Manually accept keys.
- **PasswordRuleID**: (default: 0) ID of the default password rule assigned to managed accounts created under this managed system.
- **DSSKeyRuleID**: (default: 0) ID of the default DSS key rule assigned to managed accounts created under this managed system. Can be set when **Platform.DSSFlag** is true.
- **LoginAccountID**: (optional) ID of the functional account used for SSH Session logins. Can be set if the **Platform.LoginAccountFlag** is true.
- **ReleaseDuration**: (minutes: 1-525600, default: 120) Default release duration.
- **MaxReleaseDuration**: (minutes: 1-525600, default: 525600) Default maximum release duration.
- **ISAReleaseDuration**: (minutes: 1-525600, default: 120) Default Information Systems Administrator (ISA) release duration.
- **AutoManagementFlag**: (default: false) True if password auto-management is enabled, otherwise false. Can be set if **Platform.AutoManagementFlag** is true.
  - **FunctionalAccountID**: (required if **AutoManagementFlag** is true) ID of the functional account used for local managed account password changes. **FunctionalAccount.PlatformID** must either match the **ManagedSystem.PlatformID** or be a domain platform (AD, LDAP).
  - **ElevationCommand**: (optional) Elevation command to use. Can be set if **Platform.SupportsElevationFlag** is true (sudo, pbrun, pmrun).
  - **CheckPasswordFlag**: True to enable password testing, otherwise false.
  - **ChangePasswordAfterAnyReleaseFlag**: True to change passwords on release of a request, otherwise false.
  - **ResetPasswordOnMismatchFlag**: True to queue a password change when scheduled password test fails, otherwise false.
  - **ChangeFrequencyType**: (default: first) The change frequency for scheduled password changes:
    - **first**: Changes scheduled for the first day of the month.
    - **last**: Changes scheduled for the last day of the month.
    - **xdays**: Changes scheduled every x days (ChangeFrequencyDays).
  - **ChangeFrequencyDays**: (days: 1-999, required if **ChangeFrequencyType** is **xdays**) When **ChangeFrequencyType** is **xdays**, password changes take place this configured number of days.
  - **ChangeTime**: (24hr format: 00:00-23:59, default: 23:30) UTC time of day scheduled password changes take place.
- **RemoteClientType**: The type of remote client to use.
  - **None**: No remote client.
  - **EPM**: Endpoint Privilege Management.
- **ApplicationHostID**: (default: null, required when **Platform.RequiresApplicationHost** = true) Managed system ID of the target application host. Must be an ID of a managed system whose **IsApplicationHost** = true.
- **IsApplicationHost**: (default: false) true if the managed system can be used as an application host, otherwise false. Can be set when the **Platform.ApplicationHostFlag** = true, and cannot be set when **ApplicationHostID** has a value.

#### Create Managed System using workgroup id

```
# create managed system by workgroup id / API Version: 3.0
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
  description                            = "managed_system_by_workgroup ${random_uuid.generated.result}"
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
}
# create managed system by workgroup id / API Version: 3.1
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
  description                            = "managed_system_by_workgroup ${random_uuid.generated.result}"
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
}
# create managed system by workgroup id / API Version: 3.2
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
  description                            = "managed_system_by_workgroup ${random_uuid.generated.result}"
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
# create managed system by workgroup id / API Version: 3.3
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
  description                            = "managed_system_by_workgroup ${random_uuid.generated.result}"
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
```

**Parameters Description**

See available API versions for this resource: [WorkgroupsID Managed Systems.](doc:password-safe-api#post-workgroupsidmanagedsystems)

- **WorkgroupID**: (required) Workgroup Id.
- **EntityTypeID**: (required) Type of entity being created.
- **HostName**: (required) Name of the host (applies to static asset, static database, directory, cloud). Max string length is 128 characters.
  - **Static Asset**: Asset name.
  - **Static Database**: Database host name.
  - **Directory**: Directory/domain name.
  - **Cloud**: Cloud system name.
- **IPAddress**: (required) IPv4 address of the host (applies to static asset, static database). Max string length is 45.
- **DnsName**: DNS name of the host (applies to static asset, static database). Max string length is 255.
- **InstanceName**: Name of the database instance. Required when **IsDefaultInstance** is false (applies to static database only). Max string length is 100.
- **IsDefaultInstance**: True if the database instance is the default instance, otherwise false. Only platforms MS SQL Server and MySQL support setting this value to true (applies to static database only).
- **Template**: The database connection template (applies to static database only).
- **ForestName**: Name of the directory forest (required for Active Directory; optional for Entra ID). Max string length is 64.
- **UseSSL** (default: false) True to use an SSL connection, otherwise false (applies to directory only).
- **PlatformID**: (required) ID of the managed system platform.
- **NetBiosName**: The NetBIOS name of the host. Can be set if **Platform.NetBiosNameFlag** is true. Max string length is 15.
- **ContactEmail**: Max string length is 1000.
- **Description**: Max string length is 255.
- **Port**: (optional) The port used to connect to the host. If null and the related **Platform.PortFlag** is true, Password Safe uses **Platform.DefaultPort** for communication.
- **Timeout**: (seconds, default: 30) Connection timeout. Length of time in seconds before a slow or unresponsive connection to the system fails.
- **SshKeyEnforcementMode**: (default: 0/None) Enforcement mode for SSH host keys.
  - **0**: None
  - **1**: Auto. Auto accept initial key.
  - **2**: Strict. Manually accept keys.
- **PasswordRuleID**: (default: 0) ID of the default password rule assigned to managed accounts created under this managed system.
- **DSSKeyRuleID**: (default: 0) ID of the default DSS key rule assigned to managed accounts created under this managed system. Can be set when **Platform.DSSFlag** is true.
- **LoginAccountID**: (optional) ID of the functional account used for SSH session logins. Can be set if the **Platform.LoginAccountFlag** is true.
- **AccountNameFormat**: (Active Directory only, default: 0) Account name format to use:
  - **0**: Domain and account. Use **ManagedAccount.DomainName\\ManagedAccount.AccountName**.
  - **1**: UPN. Use the managed account UPN.
  - **2**: SAM. Use the managed account SAM account name.
- **OracleInternetDirectoryID**: The Oracle Internet Directory ID (applies to database entity types and Oracle platform only).
- **OracleInternetDirectoryServiceName**: (required when **OracleInternetDirectoryID** is set) The database service name related to the given **OracleInternetDirectoryID** (applies to database entity types and Oracle platform only). Max string length is 200.
- **ReleaseDuration**: (minutes: 1-525600, default: 120) Default release duration.
- **MaxReleaseDuration**: (minutes: 1-525600, default: 525600) Default maximum release duration.
- **ISAReleaseDuration**: (minutes: 1-525600, default: 120) Default Information Systems Administrator (ISA) release duration.
- **AutoManagementFlag**: (default: false) True if password auto-management is enabled, otherwise false. Can be set if **Platform.AutoManagementFlag** is true.
  - **FunctionalAccountID**: (required if **AutoManagementFlag** is true) ID of the functional account used for local managed account password changes. **FunctionalAccount.PlatformID** must either match the **ManagedSystem.PlatformID** or be a directory platform (AD, LDAP).
  - **ElevationCommand**: (optional) Elevation command to use. Can be set if **Platform.SupportsElevationFlag** is true.
    - **sudo**
    - **pbrun**
    - **pmrun**
- **CheckPasswordFlag**: True to enable password testing, otherwise false.
- **ChangePasswordAfterAnyReleaseFlag**: True to change passwords on release of a request, otherwise false.
- **ResetPasswordOnMismatchFlag**: True to queue a password change when scheduled password test fails, otherwise false.
- **ApplicationHostID**: (default: null, required when **Platform.RequiresApplicationHost** = true) Managed system ID of the target application host. Must be an ID of a managed system where **IsApplicationHost** = true.
- **IsApplicationHost**: (default: false) true if the managed system can be used as an application host, otherwise false. Can be set when the Platform.ApplicationHostFlag = true, and cannot be set when ApplicationHostID has a value.
- **RemoteClientType**: (default: None) The type of remote client to use.
  - **None**: No remote client.
  - **EPM**: Endpoint Privilege Management.
- **AccessURL**: (default: default URL for the selected platform) The URL used for cloud access (applies to cloud systems only). Max string length is 2048.

#### Create Managed System using database id

```
# create managed system by database id
resource "passwordsafe_managed_system_by_database" "managed_system_by_database" {
  database_id                            = "2"
  contact_email                          = "admin@example.com"
  description                            = "Managed system for example DB"
  timeout                                = 30
  password_rule_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 45
  auto_management_flag                   = false
  functional_account_id                  = 0
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "xdays"
  change_frequency_days                  = 15
  change_time                            = "03:00"
}
```

**Parameters Description**

See available API versions for this resource: [Password Safe APIs](https://docs.beyondtrust.com/bips/docs/api-reference)

- **Database Id**: (required) Database Id.
- **ContactEmail**: Max string length is 1000.
- **Description**: Max string length is 255.
- **Timeout**: (seconds, default: 30) Connection timeout. Length of time in seconds before a slow or unresponsive connection to the system fails.
- **PasswordRuleID**: (default: 0) ID of the default password rule assigned to managed accounts created under this managed system.
- **ReleaseDuration**: (minutes: 1-525600, default: 120) Default release duration.
- **MaxReleaseDuration**: (minutes: 1-525600, default: 525600) Default maximum release duration.
- **ISAReleaseDuration**: (minutes: 1-525600, default: 120) Default Information Systems Administrator (ISA) release duration.
- **AutoManagementFlag**: (default: false) True if password auto-management is enabled, otherwise false. Can be set if **Platform.AutoManagementFlag** is true.
  - **FunctionalAccountID**: (required if **AutoManagementFlag** is true) ID of the functional account used for local managed account password changes. **FunctionalAccount.PlatformID** must either match the ManagedSystem.PlatformID or be a domain platform (AD, LDAP).
  - **CheckPasswordFlag**: True to enable password testing, otherwise false.
  - **ChangePasswordAfterAnyReleaseFlag**: True to change passwords on release of a request, otherwise false.
  - **ResetPasswordOnMismatchFlag**: True to queue a password change when scheduled password test fails, otherwise false.
  - **ChangeFrequencyType**: (default: first) The change frequency for scheduled password changes:
    - **first**: Changes scheduled for the first day of the month.
    - **last**: Changes scheduled for the last day of the month.
    - **xdays**: Changes scheduled every x days (see **ChangeFrequencyDays**).
  - **ChangeFrequencyDays**: (days: 1-999, required if **ChangeFrequencyType** is xdays) When **ChangeFrequencyType** is **xdays**, password changes take place this configured number of days.
  - **ChangeTime**: (24hr format: 00:00-23:59, default: 23