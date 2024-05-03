package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-freeipa/freeipa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: freeipa.Provider})
}
