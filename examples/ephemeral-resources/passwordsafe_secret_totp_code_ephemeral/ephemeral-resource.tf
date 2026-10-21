ephemeral "passwordsafe_secret_totp_code_ephemeral" "totp" {
  secret_id = passwordsafe_credential_secret.my_totp_credenial_secret.id
}
