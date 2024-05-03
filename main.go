package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	
	"github.com/rework-space-com/terraform-provider-freeipa/freeipa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider})
}
