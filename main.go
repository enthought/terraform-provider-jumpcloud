package main

import (
	"dbolack/terraform-provider-jumpcloud/jumpcloud"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return jumpcloud.Provider()
		},
	})
}
