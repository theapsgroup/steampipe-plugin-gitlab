package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"strings"
)

func tableCommit() *plugin.Table{
	return &plugin.Table{
		Name: "gitlab_commit",
		Description: "Commits in the given project.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate: listCommits,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate: getCommit,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The ID (commit hash) of the commit."},
			{Name: "short_id", Type: proto.ColumnType_STRING, Description: "The short ID (short commit hash) of the commit."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the commit."},
			{Name: "author_name", Type: proto.ColumnType_STRING, Description: "The name of the commit author."},
			{Name: "author_email", Type: proto.ColumnType_STRING, Description: "The email of the commit author."},
			{Name: "authored_date", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of commit."},
			{Name: "committer_name", Type: proto.ColumnType_STRING, Description: "The name of the committer."},
			{Name: "committer_email", Type: proto.ColumnType_STRING, Description: "The email address of the committer."},
			{Name: "committed_date", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of the commit."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of the creation of commit."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "The commit message."},
			{Name: "parent_ids", Type: proto.ColumnType_JSON, Description: "Array of parent commit hashes.", Transform: transform.FromField("ParentIDs")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project containing the commit - link to `gitlab_project.ID`"},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url of the commit.", Transform: transform.FromField("WebURL")},
		},
	}
}

func listCommits(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())
	getAll := true

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListCommitsOptions{All: &getAll, ListOptions: api.ListOptions{
		Page: 1,
		PerPage: 50,
	}}

	for {
		commits, resp, err := conn.Commits.ListCommits(projectId, opt)
		if err != nil {
			// Handle error of project id not being valid.
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			return nil, err
		}

		for _, commit := range commits {
			commit.ProjectID = projectId
			commit.Message = strings.TrimRight(commit.Message, "\n") // remove trailing newline from commit message.
			d.StreamListItem(ctx, commit)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}
	return nil, nil
}

func getCommit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())
	id := d.KeyColumnQuals["id"].GetStringValue()

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	commit, _, err := conn.Commits.GetCommit(projectId, id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	commit.ProjectID = projectId
	commit.Message = strings.TrimRight(commit.Message, "\n") // remove trailing newline from commit message.
	return commit, nil
}