---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "passwordsafe_managed_system_by_database Resource - terraform-provider-passwordsafe"
subcategory: ""
description: |-
  Managed System by Database Id Resource, creates managed system by database id.
---

# passwordsafe_managed_system_by_database (Resource)

Managed System by Database Id Resource, creates managed system by database id.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database_id` (String) Database Id

### Optional

- `auto_management_flag` (Boolean) Auto Management Flag
- `change_frequency_days` (Number) Change Frequency Days (required if ChangeFrequencyType is xdays)
- `change_frequency_type` (String) Change Frequency Type (one of: first, last, xdays)
- `change_password_after_any_release_flag` (Boolean) Change Password After Any Release Flag
- `change_time` (String) Change Time (format: HH:MM)
- `check_password_flag` (Boolean) Check Password Flag
- `contact_email` (String) Contact Email (max 1000 characters, must be a valid email)
- `description` (String) Description (max 255 characters)
- `functional_account_id` (Number) Functional Account ID (required if AutoManagementFlag is true)
- `isa_release_duration` (Number) ISA Release Duration (min: 1, max: 525600)
- `max_release_duration` (Number) Max Release Duration (min: 1, max: 525600)
- `password_rule_id` (Number) Password Rule ID
- `release_duration` (Number) Release Duration (min: 1, max: 525600)
- `reset_password_on_mismatch_flag` (Boolean) Reset Password On Mismatch Flag
- `timeout` (Number) Timeout (min: 1)

### Read-Only

- `managed_system_id` (Number) Managed System Id
- `managed_system_name` (String) Managed System Name
