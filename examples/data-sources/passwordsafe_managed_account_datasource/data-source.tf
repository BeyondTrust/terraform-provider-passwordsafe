// managed_accounts_list
data "passwordsafe_managed_account_datasource" "managed_accounts_list" {
}
output "managed_accounts_list" {
  value = data.passwordsafe_managed_account_datasource.managed_accounts_list.managed_accounts[0].account_name
}