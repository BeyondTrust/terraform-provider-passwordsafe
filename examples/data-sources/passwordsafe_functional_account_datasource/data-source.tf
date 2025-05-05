// functional_accounts_list
data "passwordsafe_functional_account_datasource" "functional_accounts_list" {
}
output "functional_account_list" {
  value = [
    for acc in data.passwordsafe_functional_account_datasource.functional_accounts_list.accounts : acc.account_name
    if acc.account_name == "svc-monitoring"
  ][0]
}