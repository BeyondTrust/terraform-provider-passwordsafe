# create managed system by database id
resource "passwordsafe_managed_system_by_database" "managed_system_by_database" {
  database_id                            = "2"
  contact_email                          = "admin@example.com"
  description                            = "Managed system for example DB"
  timeout                                = 30
  password_rule_id                       = 0
  release_duration                       = 60
  max_release_duration                   = 120
  isa_release_duration                   = 45
  auto_management_flag                   = false
  functional_account_id                  = 0
  check_password_flag                    = true
  change_password_after_any_release_flag = false
  reset_password_on_mismatch_flag        = true
  change_frequency_type                  = "xdays"
  change_frequency_days                  = 15
  change_time                            = "03:00"
}