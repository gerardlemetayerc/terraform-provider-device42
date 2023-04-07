package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/monprojet/mon-provider-terraform-device42/device42"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: device42.Provider,
	})
}
