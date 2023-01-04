package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableSnippet() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_snippet",
		Description: "The current logged in users GitLab snippets",
		List: &plugin.ListConfig{
			Hydrate: listSnippets,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the snippet."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the snippet."},
			{Name: "file_name", Type: proto.ColumnType_STRING, Description: "The file name of the snippet."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the snippet."},
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "The ID of the author - link to `gitlab_user.id`", Transform: transform.FromField("Author.ID")},
			{Name: "author_username", Type: proto.ColumnType_STRING, Description: "The username of the author -  - link to `gitlab_user.username`", Transform: transform.FromField("Author.Username")},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The display name of the author.", Transform: transform.FromField("Author.Name")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of the creation of the snippet."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp that the snippet was last updated."},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url to the snippet.", Transform: transform.FromField("WebURL")},
			{Name: "raw_url", Type: proto.ColumnType_STRING, Description: "The url to the raw content of the snippet.", Transform: transform.FromField("RawURL")},
			{Name: "files", Type: proto.ColumnType_JSON, Description: "An array of file paths & urls."},
		},
	}
}

func listSnippets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListSnippetsOptions{Page: 1, PerPage: 30}

	for {
		snippets, resp, err := conn.Snippets.ListSnippets(opt)
		if err != nil {
			return nil, err
		}

		for _, snippet := range snippets {
			d.StreamListItem(ctx, snippet)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
