package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectDeployToken() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_deploy_token",
		Description: "Obtain information about deploy tokens for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectDeployTokens,
		},
		Columns: projectDeployTokenColumns(),
	}
}

// Hydrate Functions
func listProjectDeployTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectDeployTokens", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectDeployTokens", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectDeployTokensOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectDeployTokens", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		tokens, resp, err := conn.DeployTokens.ListProjectDeployTokens(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectDeployTokens", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain deploy tokens for project_id %d\n%v", projectId, err)
		}

		for _, token := range tokens {
			d.StreamListItem(ctx, token)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectDeployTokens", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectDeployTokens", "completed successfully")
	return nil, nil
}

// Column Function
func projectDeployTokenColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the deploy token.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the deploy token.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "The username to use with the deploy token.",
		},
		{
			Name:        "scopes",
			Type:        proto.ColumnType_JSON,
			Description: "The scopes granted to the deploy token.",
		},
		{
			Name:        "expires_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date of expiration of the deploy token.",
		},
		{
			Name:        "revoked",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the deploy token has been revoked.",
			Transform:   transform.FromField("Revoked"),
		},
		{
			Name:        "expired",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the deploy token has expired.",
			Transform:   transform.FromField("Expired"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
