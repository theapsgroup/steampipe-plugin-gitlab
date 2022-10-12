package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"steampipe-plugin-gitlab/gitlab"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: gitlab.Plugin})
}
