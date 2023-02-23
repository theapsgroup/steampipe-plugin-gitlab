package gitlab

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_version",
		Description: "Obtain information about the version of the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listVersion,
		},
		Columns: versionColumns(),
	}
}

// Hydrate Function
func listVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listVersion", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listVersion", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	versionData, _, err := conn.Version.GetVersion()
	if err != nil {
		plugin.Logger(ctx).Error("listVersion", "error", err)
		return nil, fmt.Errorf("unable to obtain version information\n%v", err)
	}

	d.StreamListItem(ctx, versionData)

	plugin.Logger(ctx).Debug("listVersion", "completed successfully")
	return nil, nil
}

// Column Function
func versionColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "version",
			Type:        proto.ColumnType_STRING,
			Description: "GitLab Version",
		},
		{
			Name:        "revision",
			Type:        proto.ColumnType_STRING,
			Description: "Revision of the current version",
		},
	}
}
