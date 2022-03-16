package main

import (
	"github.com/enthought/terraform-provider-jumpcloud/jumpcloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {

	var debugMode bool

	opts := &plugin.ServeOpts{
		Debug: debugMode,

		ProviderAddr: "terraform.enthought.com/enthought/jumpcloud",

		ProviderFunc: jumpcloud.New(version),
	}

	plugin.Serve(opts)
}
