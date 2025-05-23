---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "passwordsafe_database_datasource Data Source - terraform-provider-passwordsafe"
subcategory: ""
description: |-
  Database Datasource, get databases list.
---

# passwordsafe_database_datasource (Data Source)

Database Datasource, get databases list.

## Example Usage

```terraform
// databases_list
data "passwordsafe_database_datasource" "databases_list" {
}
output "database_list" {
  value = data.passwordsafe_database_datasource.databases_list.databases[1].instance_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `databases` (Block List) Database Datasource Attributes (see [below for nested schema](#nestedblock--databases))

<a id="nestedblock--databases"></a>
### Nested Schema for `databases`

Read-Only:

- `asset_id` (Number) Asset ID
- `database_id` (Number) Database ID
- `instance_name` (String) Instance Name
- `is_default_instance` (Boolean) Is Default Instance
- `platform_id` (Number) Platform ID
- `port` (Number) Port
- `template` (String) Template
- `version` (String) Version
