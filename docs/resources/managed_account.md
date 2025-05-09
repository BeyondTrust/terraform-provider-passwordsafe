---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "passwordsafe_managed_account Resource - terraform-provider-passwordsafe"
subcategory: ""
description: |-
  Managed Account Resource, creates managed account.
---

# passwordsafe_managed_account (Resource)

Managed Account Resource, creates managed account.

## Example Usage

```terraform
resource "passwordsafe_managed_account" "my_managed_account" {
  system_name  = "system_integration_test"
  account_name = "managed_account_${random_uuid.generated.result}"
  password     = "MyTest101*!"
  api_enabled  = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_name` (String)
- `password` (String)
- `system_name` (String)

### Optional

- `api_enabled` (Boolean)
- `auto_management_flag` (Boolean)
- `change_com_plus_flag` (Boolean)
- `change_dcom_flag` (Boolean)
- `change_frequency_days` (Number)
- `change_frequency_type` (String)
- `change_password_after_any_release_flag` (Boolean)
- `change_scom_flag` (Boolean)
- `change_services_flag` (Boolean)
- `change_tasks_flag` (Boolean)
- `change_time` (String)
- `change_windows_auto_logon_flag` (Boolean)
- `check_password_flag` (Boolean)
- `description` (String)
- `distinguished_name` (String)
- `domain_name` (String)
- `dss_auto_management_flag` (Boolean)
- `isa_release_duration` (Number)
- `login_account_flag` (Boolean)
- `max_concurrent_requests` (Number)
- `max_release_duration` (Number)
- `next_change_date` (String)
- `object_id` (String)
- `passphrase` (String)
- `password_fallback_flag` (Boolean)
- `password_rule_id` (Number)
- `private_key` (String)
- `release_duration` (Number)
- `release_notification_email` (String)
- `reset_password_on_mismatch_flag` (Boolean)
- `restart_services_flag` (Boolean)
- `sam_account_name` (String)
- `use_own_credentials` (Boolean)
- `user_principal_name` (String)
- `workgroup_id` (Number)

### Read-Only

- `id` (String) The ID of this resource.
