resource "passwordsafe_database" "database" {
  asset_id            = "1"
  platform_id         = 10
  instance_name       = "primary-db-instance"
  is_default_instance = false
  port                = 5432
  version             = "13.3"
  template            = "standard-template"
}
