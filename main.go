package main

import (
	"github.com/theapsgroup/steampipe-plugin-gitlab/gitlab"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: gitlab.Plugin})
}
