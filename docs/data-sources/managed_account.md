---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "passwordsafe_managed_account Data Source - terraform-provider-passwordsafe"
subcategory: ""
description: |-
  Managed Account Datasource, gets managed account.
---

# passwordsafe_managed_account (Data Source)

Managed Account Datasource, gets managed account.

## Example Usage

```terraform
data "passwordsafe_managed_account" "manage_account_01" {
  system_name  = "system01"
  account_name = "managed_account02"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_name` (String)
- `system_name` (String)

### Optional

- `value` (String)

### Read-Only

- `id` (String) The ID of this resource.
