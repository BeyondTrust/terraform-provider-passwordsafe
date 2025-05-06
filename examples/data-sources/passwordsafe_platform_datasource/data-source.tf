// plarforms_list
data "passwordsafe_platform_datasource" "plarforms_list" {
}
output "plarforms_list" {
  value = data.passwordsafe_platform_datasource.plarforms_list.platforms[0].name
}
