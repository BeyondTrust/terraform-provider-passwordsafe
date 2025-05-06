// folders_list
data "passwordsafe_folder_datasource" "folders_list" {
}
output "folders_list" {
  value = data.passwordsafe_folder_datasource.folders_list.folders[0].name
}