// safes_list
data "passwordsafe_safe_datasource" "safes_list" {
}
output "safes_list" {
  value = data.passwordsafe_safe_datasource.safes_list.safes[0].name
}