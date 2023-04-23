package main

import (
	device42 "github.com/gerardlemetayerc/terraform-provider-device42/Device42"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: device42.Provider,
	})
}
