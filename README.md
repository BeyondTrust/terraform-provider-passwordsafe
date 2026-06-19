<a href="https://www.beyondtrust.com">
    <img src="images/beyondtrust_logo.svg" alt="BeyondTrust" title="BeyondTrust" align="right" height="50">
</a>

# BeyondTrust Password Safe Terraform Provider

[![Tests](https://github.com/BeyondTrust/terraform-provider-passwordsafe/actions/workflows/tests.yml/badge.svg)](https://github.com/BeyondTrust/terraform-provider-passwordsafe/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/BeyondTrust/terraform-provider-passwordsafe)](https://goreportcard.com/report/github.com/BeyondTrust/terraform-provider-passwordsafe)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/BeyondTrust/terraform-provider-passwordsafe/badge)](https://scorecard.dev/viewer/?uri=github.com/BeyondTrust/terraform-provider-passwordsafe)
[![Release](https://github.com/BeyondTrust/terraform-provider-passwordsafe/actions/workflows/release.yml/badge.svg)](https://github.com/BeyondTrust/terraform-provider-passwordsafe/actions/workflows/release.yml)
[![Registry](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fregistry.terraform.io%2Fv1%2Fproviders%2FBeyondTrust%2Fpasswordsafe&query=%24.version&label=registry&color=7B42BC&logo=terraform)](https://registry.terraform.io/providers/BeyondTrust/passwordsafe/latest)
[![License](https://img.shields.io/github/license/BeyondTrust/terraform-provider-passwordsafe)](LICENSE)

The [BeyondTrust Password Safe Terraform Provider](https://registry.terraform.io/providers/BeyondTrust/passwordsafe/latest/docs) allows [Terraform](https://terraform.io) to manage access to resources in Password Safe.  Terraform configuration files can be configured to retrieve secrets from Password Safe and permissions for access to secrets in Password Safe can be granted to specific accounts within BeyondInsight.

## Build Requirements
If you wish to build this project locally, please complete the following steps.

1. Install Go.  Instructions are on [their website](https://go.dev/doc/install).
1. Install Terraform.  Instructions are on [their website](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli).

1. Clone this repository.

    ```bash
    git clone https://github.com/BeyondTrust/terraform-provider-passwordsafe
    ```

1. Generate the provider binary file from passwordsafe-integrations-terraform folder.

    ```bash
    go build -o terraform-provider-passwordsafe_major_minor_build
    ```


    To run unit tests you can use:

    ```bash
   go test ./...
    ```

1. Move the binary file to the proper directory to be recognized by Terraform.

1. Update _terraform/main.tf_ and _terraform/terraform.tfvars_ files according to your needs, and then run terraform commands.

    ```bash
    terraform init
    terraform plan
    ```

## Usage
The [product documentation](https://www.beyondtrust.com/docs/beyondinsight-password-safe/) as well as the Terraform-specific [usage documentation](https://www.beyondtrust.com/docs/beyondinsight-password-safe/ps/integrations/terraform/index.htm) are hosted on our website.

## License
This software is distributed under the GNU General Public License (GPL) v3.0 License. See `LICENSE.txt` for more information.

## Get Help
Contact [BeyondTrust support](https://www.beyondtrust.com/docs/index.htm#support)