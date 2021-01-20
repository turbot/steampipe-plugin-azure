package main

import (
	"github.com/turbot/steampipe-plugin-azure/azure"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: azure.Plugin})
}
