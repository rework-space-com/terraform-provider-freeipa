resource freeipa_user "user-0" {
    name = "test-user2"
    first_name = "Test"
    last_name = "User"
    city = "Brussels"
    login_shell = "/bin/bash"
    home_directory = "/opt/users/test-user"
}

# data freeipa_user "user-0" {
#     name = freeipa_user.user-0.name
# }

# output "user-0" {
#     value = data.freeipa_user.user-0
# }