package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectDeployKey() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_deploy_key",
		Description: "Obtain information about deploy keys for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectDeployKeys,
		},
		Columns: projectDeployKeyColumns(),
	}
}

// Hydrate Functions
func listProjectDeployKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectDeployKeys", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectDeployKeys", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectDeployKeysOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectDeployKeys", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		keys, resp, err := conn.DeployKeys.ListProjectDeployKeys(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectDeployKeys", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain deploy keys for project_id %d\n%v", projectId, err)
		}

		for _, key := range keys {
			d.StreamListItem(ctx, key)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectDeployKeys", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectDeployKeys", "completed successfully")
	return nil, nil
}

// Column Function
func projectDeployKeyColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the deploy key.",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the deploy key.",
		},
		{
			Name:        "key",
			Type:        proto.ColumnType_STRING,
			Description: "The public SSH key material of the deploy key.",
		},
		{
			Name:        "can_push",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the deploy key is allowed to push to the repository.",
			Transform:   transform.FromField("CanPush"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the deploy key was created.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
