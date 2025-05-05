resource "passwordsafe_workgroup" "workgroup" {
  name = "workgroup_name_${random_uuid.generated.result}"
}