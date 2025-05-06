resource "passwordsafe_text_secret" "my_text_secret" {
  folder_name = "folder1"
  title       = "Text_Secret"
  description = "my text secret description"
  owner_type  = "User"
  text        = "password_text"
  notes       = "My notes"
  group_id    = 1
}
