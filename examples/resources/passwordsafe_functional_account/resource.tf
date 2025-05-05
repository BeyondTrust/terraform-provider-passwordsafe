resource "passwordsafe_functional_account" "functional_account" {
  platform_id           = 1
  domain_name           = "test.example.com"
  account_name          = "FUNCTIONAL_ACCOUNT"
  display_name          = "FUNCTIONAL_ACCOUNT"
  password              = "pass-value"
  private_key           = "private key value"
  passphrase            = "my-passphrase"
  description           = "functional account description"
  elevation_command     = "sudo"
  tenant_id             = ""
  object_id             = ""
  secret                = "super-secret-value"
  service_account_email = "test@test.com"
  azure_instance        = "AzurePublic"
}