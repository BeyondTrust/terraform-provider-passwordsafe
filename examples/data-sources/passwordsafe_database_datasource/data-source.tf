// databases_list
data "passwordsafe_database_datasource" "databases_list" {
}
output "database_list" {
  value = data.passwordsafe_database_datasource.databases_list.databases[1].instance_name
}
