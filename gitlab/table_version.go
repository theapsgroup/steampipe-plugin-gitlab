package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_version",
		Description: "GitLab version information",
		List: &plugin.ListConfig{
			Hydrate: listVersion,
		},
		Columns: []*plugin.Column{
			{Name: "version", Type: proto.ColumnType_STRING, Description: "GitLab Version"},
			{Name: "revision", Type: proto.ColumnType_STRING, Description: "Revision of the current version"},
		},
	}
}

func listVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	versionData, _, err := conn.Version.GetVersion()
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, versionData)
	return nil, nil
}
