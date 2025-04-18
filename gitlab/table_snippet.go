package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableSnippet() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_snippet",
		Description: "Obtain information about snippets for the currently authenticated user within the GitLab instance.",
		List: &plugin.ListConfig{
			Hydrate: listSnippets,
		},
		Columns: snippetColumns(),
	}
}

func listSnippets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listSnippets", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listSnippets", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	opt := &api.ListSnippetsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listSnippets", "page", opt.Page, "perPage", opt.PerPage)
		snippets, resp, err := conn.Snippets.ListSnippets(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listSnippets", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain snippets for authenticated user\n%v", err)
		}

		for _, snippet := range snippets {
			d.StreamListItem(ctx, snippet)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listSnippets", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listSnippets", "completed successfully")
	return nil, nil
}

// Column Function
func snippetColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the snippet.",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the snippet.",
		},
		{
			Name:        "file_name",
			Type:        proto.ColumnType_STRING,
			Description: "The file name of the snippet.",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "The description of the snippet.",
		},
		{
			Name:        "author_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the author - link to `gitlab_user.id`",
			Transform:   transform.FromField("Author.ID"),
		},
		{
			Name:        "author_username",
			Type:        proto.ColumnType_STRING,
			Description: "The username of the author -  - link to `gitlab_user.username`",
			Transform:   transform.FromField("Author.Username"),
		},
		{
			Name:        "author_name",
			Type:        proto.ColumnType_STRING,
			Description: "The display name of the author.",
			Transform:   transform.FromField("Author.Name"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of the creation of the snippet.",
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp that the snippet was last updated.",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to the snippet.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "raw_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url to the raw content of the snippet.",
			Transform:   transform.FromField("RawURL"),
		},
		{
			Name:        "files",
			Type:        proto.ColumnType_JSON,
			Description: "An array of file paths & urls.",
		},
	}
}
