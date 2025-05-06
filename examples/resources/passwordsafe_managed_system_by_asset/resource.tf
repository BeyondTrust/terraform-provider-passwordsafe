# create managed system by asset id / API Version: 3.0
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
  port                                   = 5432
  timeout                                = 30
  ssh_key_enforcement_mode               = 1
  password_rule_id                       = 0
  dss_key_rule_id                        = 0
  login_account_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 90
  auto_management_flag                   = false
  functional_account_id                  = 0
  elevation_command                      = "sudo su"
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "last"
  change_frequency_days                  = 30
  change_time                            = "02:00"
}

# create managed system by asset id / API Version: 3.1
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
  port                                   = 5432
  timeout                                = 30
  ssh_key_enforcement_mode               = 1
  password_rule_id                       = 0
  dss_key_rule_id                        = 0
  login_account_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 90
  auto_management_flag                   = false
  functional_account_id                  = 0
  elevation_command                      = "sudo su"
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "last"
  change_frequency_days                  = 30
  change_time                            = "02:00"
  remote_client_type                     = "EPM"
}

# create managed system by asset id / API Version: 3.2
resource "passwordsafe_managed_system_by_asset" "managed_system_by_asset" {
  asset_id                               = "48"
  platform_id                            = 2
  contact_email                          = "admin@example.com"
  description                            = "managed_system_by_asset"
  port                                   = 5432
  timeout                                = 30
  ssh_key_enforcement_mode               = 1
  password_rule_id                       = 0
  dss_key_rule_id                        = 0
  login_account_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 90
  auto_management_flag                   = false
  functional_account_id                  = 0
  elevation_command                      = "sudo su"
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "last"
  change_frequency_days                  = 30
  change_time                            = "02:00"
  remote_client_type                     = "EPM"
  application_host_id                    = 0
  is_application_host                    = false
}
