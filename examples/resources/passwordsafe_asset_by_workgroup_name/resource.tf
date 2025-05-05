resource "passwordsafe_asset_by_workgroup_name" "asset_by_workgroup_name" {
  work_group_name  = passwordsafe_workgroup.workgroup.name
  ip_address       = "10.20.30.40"
  asset_name       = "Prod_Server_03"
  dns_name         = "server01.company.com"
  domain_name      = "company.com"
  mac_address      = "00:1A:2B:3C:4D:5E"
  asset_type       = "Windows Server"
  description      = "Production Windows Server hosting critical applications"
  operating_system = "Windows Server 2022"
}