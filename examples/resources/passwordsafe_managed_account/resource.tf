resource "passwordsafe_managed_account" "my_managed_account" {
  system_name  = "system_integration_test"
  account_name = "managed_account_${random_uuid.generated.result}"
  password     = "MyTest101*!"
  api_enabled  = true
}