// current TOTP code of a credential secret (API version 3.3)
data "passwordsafe_secret_totp_code" "totp" {
  secret_id = passwordsafe_credential_secret.my_totp_credenial_secret.id
}
output "totp_code" {
  value     = data.passwordsafe_secret_totp_code.totp.code
  sensitive = true
}
