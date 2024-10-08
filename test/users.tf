resource freeipa_user "user-0" {
    name = "test-user0"
    first_name = "Test"
    last_name = "User"
    city = "Brussels"
    login_shell = "/bin/bash"
    home_directory = "/opt/users/test-user0"
}

resource freeipa_user "user-1" {
    name = "test-user1"
    first_name = "Test"
    last_name = "User"
    city = "Brussels"
    login_shell = "/bin/bash"
    home_directory = "/opt/users/test-user-1"
    employee_type = "Patsy"
    account_disabled = false
}
# data freeipa_user "user-0" {
#     name = freeipa_user.user-0.name
# }

# output "user-0" {
#     value = data.freeipa_user.user-0
# }