package main

import (
	"github.com/colinhoglund/terraform-provider-kops/kops"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kops.Provider})
}
