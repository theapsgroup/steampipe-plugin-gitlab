package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectAccessRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_access_request",
		Description: "Obtain access requests for a specific project in the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listProjectAccessRequests,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: projectAccessRequestColumns(),
	}
}

// Hydrate Function
func listProjectAccessRequests(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectAccessRequests", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectAccessRequests", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListAccessRequestsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectAccessRequests", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		reqs, resp, err := conn.AccessRequests.ListProjectAccessRequests(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectAccessRequests", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain access requests for project_id %d\n%v", projectId, err)
		}

		for _, req := range reqs {
			d.StreamListItem(ctx, req)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectAccessRequests", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectAccessRequests", "completed successfully")
	return nil, nil
}

// Column Function
func projectAccessRequestColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the access request.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the user requesting access.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the user requesting access.",
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the access request.",
		},
		{
			Name:        "access_level",
			Type:        proto.ColumnType_INT,
			Description: "The numeric value of the access level requested by the user.",
		},
		{
			Name:        "access_level_description",
			Type:        proto.ColumnType_STRING,
			Description: "The access level requested by the user.",
			Transform:   transform.FromField("AccessLevel").Transform(accessLevelTransform),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of access request creation.",
		},
		{
			Name:        "requested_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of access request submission.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project - link to `gitlab_project.id",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
