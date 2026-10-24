package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectAccessToken() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_access_token",
		Description: "Obtain information about access tokens for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectAccessTokens,
		},
		Columns: projectAccessTokenColumns(),
	}
}

// Hydrate Functions
func listProjectAccessTokens(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectAccessTokens", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectAccessTokens", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectAccessTokensOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectAccessTokens", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		tokens, resp, err := conn.ProjectAccessTokens.ListProjectAccessTokens(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectAccessTokens", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain access tokens for project_id %d\n%v", projectId, err)
		}

		for _, token := range tokens {
			d.StreamListItem(ctx, token)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectAccessTokens", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectAccessTokens", "completed successfully")
	return nil, nil
}

// Column Function
func projectAccessTokenColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the access token.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the access token.",
		},
		{
			Name:        "user_id",
			Type:        proto.ColumnType_INT,
			Description: "The id of the service account user the token authenticates as - link to `gitlab_user.id`.",
		},
		{
			Name:        "scopes",
			Type:        proto.ColumnType_JSON,
			Description: "The scopes granted to the access token.",
		},
		{
			Name:        "active",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the access token is active.",
			Transform:   transform.FromField("Active"),
		},
		{
			Name:        "revoked",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the access token has been revoked.",
			Transform:   transform.FromField("Revoked"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the access token was created.",
		},
		{
			Name:        "last_used_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of when the access token was last used.",
		},
		{
			Name:        "expires_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date of expiration of the access token.",
			Transform:   transform.FromField("ExpiresAt").NullIfZero().Transform(isoTimeTransform),
		},
		{
			Name:        "access_level",
			Type:        proto.ColumnType_INT,
			Description: "The access level granted to the access token.",
			Transform:   transform.FromField("AccessLevel"),
		},
		{
			Name:        "access_level_description",
			Type:        proto.ColumnType_STRING,
			Description: "The descriptive of the access level granted to the access token.",
			Transform:   transform.FromField("AccessLevel").Transform(accessLevelTransform),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
