## Developing the provider

Thank you for your interest in contributing.

## Documentation

Terraform [provider development documentation](https://www.terraform.io/docs/extend/) provides a good start into developing an understanding of provider development.


## Building the provider

There is a [makefile](../GNUmakefile) to help build the provider. You can build the provider by running `make build`.

```shell
$ make build
```

### Tests
To run the full suite of Acceptance tests export the required environment variables and run the `make testacc`.


```shell
$ export FREEIPA_HOST=ipa.ipatest.lan
$ export FREEIPA_USERNAME=admin
$ export FREEIPA_PASSWORD=P@ssword
$ make testacc
```

## Install provider locally
You can install provider locally for development and testing.
Use the `make install` command to compile the provider into a binary and install it in your `GOBIN` path.
```shell
$ make install
```

Terraform allows you to use local provider builds by setting a `dev_overrides` block in a configuration file called `.terraformrc`. This block overrides all other configured installation methods.

Create a new file called `.terraformrc` in your home directory (`~`), then add the `dev_overrides` block below. Change the `<PATH>` to the value returned from the `go env GOBIN` command.

If the `GOBIN` go environment variable is not set, use the default path, `/home/<Username>/go/bin`.

```terraform
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/freeipa" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}

```

### Locally installed provider configuration
Use the folloving provider configration block for the locally installed provider
```terraform
terraform {
  required_providers {
    freeipa = {
      source = "hashicorp.com/edu/freeipa"
    }
  }
}

provider "freeipa" {
  host = "ipa.ipatest.lan"
  username = "admin"
  password = "P@ssword"
  insecure = true
}
```

### Create documentation

When creating or updating resources/data resources please make sure to update the examples in the respective folder (`./examples/resources/<name>` for resources, `./examples/data-sources/<name>` for data sources)

Next you can use the following command to generate the terraform documentation from go files

```shell
make doc
```

### Start FreeIPA locally
You can start local FreeIPA instance using the provided Docker Compose file.

Add a record for the FreeIPA test server to the `/etc/hosts`.
```shell
$ sudo echo "127.0.0.1 ipa.ipatest.lan" | sudo tee -a /etc/hosts
```

In another terminal window, navigate to the `docker_compose` directory.
```bash
$ cd docker_compose
```

Run `docker-compose up` to spin up a local instance of FreeIPA.
```bash
$ docker-compose up
```

Leave this process running in your terminal window. Wait until the FreeIPA server configuration is finished. In the original terminal window, verify that the FreeIPA server is running by sending the following requests:

```bash
$ curl --insecure -s -c ./cookie.txt -k -H "Content-Type: application/x-www-form-urlencoded" -H "Accept: text/plain" -d "user=admin&password=P@ssword" https://ipa.ipatest.lan/ipa/session/login_password
$ curl --insecure  -k -b ./cookie.txt -H "Referer: https://ipa.ipatest.lan/ipa" -H "Content-Type:application/json" https://ipa.ipatest.lan/ipa/json --data '{"method":"env","params":[["version"],{}],"id":0}'
```
You should get a response from the FreeIPA API server that contains the server version and information about the principal used for authentication:
```shell
{"result": {"result": {"version": "4.10.1"}, "count": 1, "total": 120, "summary": null, "messages": [{"type": "warning", "name": "VersionMissing", "message": "API Version number was not sent, forward compatibility not guaranteed. Assuming server's API version, 2.251", "code": 13001, "data": {"server_version": "2.251"}}]}, "error": null, "id": 0, "principal": "admin@IPATEST.LAN", "version": "4.10.1"}
```
### Run individual acceptance test
Once the environment is set, you can run a specific test with this command (example for `TestAccFreeIPASudoRuleAllowCmdMembership_CaseInsensitive`)

```shell
TF_ACC=1 /usr/local/go/bin/go test -timeout 30s -run ^TestAccFreeIPASudoRuleAllowCmdMembership_CaseInsensitive$ github.com/rework-space-com/terraform-provider-freeipa/freeipa
```
