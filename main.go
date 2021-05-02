package main

import (
	"github.com/fossilmfc/terraform-provider-websupportsk/websupportsk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: websupportsk.Provider,
	})
}
