package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableApplication() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_application",
		Description: "Obtain information about OAuth applications within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listApplications,
		},
		Columns: applicationColumns(),
	}
}

// Hydrate Functions
func listApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listApplications", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listApplications", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	opt := &api.ListApplicationsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listApplications", "page", opt.Page, "perPage", opt.PerPage)
		apps, resp, err := conn.Applications.ListApplications(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listApplications", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain oauth applications\n%v", err)
		}

		for _, app := range apps {
			d.StreamListItem(ctx, app)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listApplications", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listApplications", "completed successfully")
	return nil, nil
}

// Column Function
func applicationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the application.",
		},
		{
			Name:        "application_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique identifier of the application.",
		},
		{
			Name:        "application_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the application.",
		},
		{
			Name:        "callback_url",
			Type:        proto.ColumnType_STRING,
			Description: "The redirect/callback url of the application.",
			Transform:   transform.FromField("CallbackURL"),
		},
		{
			Name:        "confidential",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the application is confidential.",
		},
	}
}
