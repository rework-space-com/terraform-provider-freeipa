resource "freeipa_user" "user-1" {
  first_name        = "Roman"
  last_name         = "Roman"
  name              = "roman"
  telephone_numbers = ["+380982555429", "2-10-11"]
  email_address     = ["roman@example.com"]
}

resource "freeipa_user" "user-2" {
  first_name        = "Roman"
  last_name         = "Roman"
  name              = "roman"
  state             = "staged"
  telephone_numbers = ["+380982555429", "2-10-11"]
  email_address     = ["roman@example.com"]
}

# Requires Terraform 1.11 or later. The password is never stored in the state.
# Increment password_wo_version to push a new password to FreeIPA.
resource "freeipa_user" "user-3" {
  first_name          = "Jane"
  last_name           = "Doe"
  name                = "jdoe"
  password_wo         = var.jdoe_password
  password_wo_version = 1
}
