resource "passwordsafe_asset_by_workgroup_id" "asset_by_workgroup_id" {
  work_group_id    = "28"
  ip_address       = "10.20.30.40"
  asset_name       = "Prod_Server_03_${random_uuid.generated.result}"
  dns_name         = "server01.company.com"
  domain_name      = "company.com"
  mac_address      = "00:1A:2B:3C:4D:5E"
  asset_type       = "Windows Server"
  description      = "Production Windows Server hosting critical applications"
  operating_system = "Windows Server 2022"
}