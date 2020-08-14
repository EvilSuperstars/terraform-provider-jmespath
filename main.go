package main

import (
	"github.com/EvilSuperstars/terraform-provider-jmespath/jmespath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jmespath.Provider,
	})
}
