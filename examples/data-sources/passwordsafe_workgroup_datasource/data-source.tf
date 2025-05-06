// workgroups_list
data "passwordsafe_workgroup_datasource" "workgroups_list" {
}
output "workgroups_list" {
  value = data.passwordsafe_workgroup_datasource.workgroups_list.workgroups[0].name
}