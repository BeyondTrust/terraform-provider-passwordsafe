// managed_system_list
data "passwordsafe_managed_system_datasource" "managed_system_list" {
}
output "managed_system_list" {
  value = data.passwordsafe_managed_system_datasource.managed_system_list.managed_systems[0].system_name
}