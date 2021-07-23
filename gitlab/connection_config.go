package gitlab

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type GitLabConfig struct {
	BaseUrl *string `cty:"baseurl"`
	Token   *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"baseurl": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &GitLabConfig{}
}

func GetConfig(connection *plugin.Connection) GitLabConfig {
	if connection == nil || connection.Config == nil {
		return GitLabConfig{}
	}

	config, _ := connection.Config.(GitLabConfig)
	return config
}
