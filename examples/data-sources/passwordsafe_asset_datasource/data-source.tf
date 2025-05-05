// assets_list
data "passwordsafe_asset_datasource" "assets_list" {
  parameter = "1" // could be workgroup id / workgroup name.
}
output "assets_list" {
  value = data.passwordsafe_asset_datasource.assets_list.assets[0].asset_name

}