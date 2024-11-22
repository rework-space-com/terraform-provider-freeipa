Terraform FreeIPA Provider
============================
Tested on FreeIPA version 4.9.1  
Download provider from [registry.terraform.io](https://registry.terraform.io/providers/rework-space-com/freeipa/latest)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 1.0+
- [Go](https://golang.org/doc/install) 1.22+ (to build the provider plugin)

Building The Provider
---------------------

Clone the repository. Enter the provider directory and build the provider

```sh
$ cd terraform-provider-freeipa
$ go build -o ~/go/bin/terraform-provider-freeipa
```
## Contributing to the provider

To contribute, please read the [contribution guidelines](_about/CONTRIBUTING.md). You may also [report an issue](https://github.com/rework-space-com/terraform-provider-freeipa/issues/new/choose).