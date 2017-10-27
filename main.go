package main

import (
	"github.com/ewbankkit/terraform-provider-jmespath/jmespath"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jmespath.Provider,
	})
}
