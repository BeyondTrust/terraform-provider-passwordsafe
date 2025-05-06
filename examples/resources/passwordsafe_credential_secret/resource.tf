resource "passwordsafe_credential_secret" "my_credenial_secret" {
  folder_name = "folder1"
  title       = "Credential_Secret"
  description = "my credential secret description"
  username    = "my_user_name"
  password    = "password_content"
  owner_type  = "User"
  notes       = "My Notes"
  group_id    = 1
}