resource "passwordsafe_file_secret" "my_file_secret" {
  folder_name  = "folder1"
  title        = "File_Secret"
  description  = "my file secret description"
  owner_type   = "User"
  file_content = file("test_secret.txt")
  file_name    = "my_secret.txt"
  notes        = "My notes"
  group_id     = 1
}
