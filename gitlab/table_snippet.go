package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			{Name: "id", Type: proto.ColumnType_INT, Description: ""},
			{Name: "title", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "file_name", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "description", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "author_id", Type: proto.ColumnType_INT, Description: "", Transform: transform.FromField("Author.ID")},
			{Name: "author_username", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("Author.Username")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: ""},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("WebURL")},
			{Name: "raw_url", Type: proto.ColumnType_STRING, Description: "", Transform: transform.FromField("RawURL")},
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

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
