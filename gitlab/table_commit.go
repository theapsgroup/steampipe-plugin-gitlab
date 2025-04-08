package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableCommit() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_commit",
		Description: "Obtain information about commits for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listCommits,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getCommit,
		},
		Columns: commitColumns(),
	}
}

// Hydrate Functions
func listCommits(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listCommits", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCommits", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	bTrue := true
	opt := &api.ListCommitsOptions{All: &bTrue, WithStats: &bTrue, ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listCommits", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		commits, resp, err := conn.Commits.ListCommits(projectId, opt)
		if err != nil {
			// Handle error of project id not being valid.
			if strings.Contains(err.Error(), "404") {
				plugin.Logger(ctx).Warn("listCommits", "projectId", projectId, "no project was found, returning empty result set")
				return nil, nil
			}
			plugin.Logger(ctx).Error("listCommits", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain commits for project_id %d\n%v", projectId, err)
		}

		for _, commit := range commits {
			commit.ProjectID = projectId
			commit.Message = strings.TrimRight(commit.Message, "\n") // remove trailing newline from commit message.
			d.StreamListItem(ctx, commit)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listCommits", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listCommits", "completed successfully")
	return nil, nil
}

func getCommit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getCommit", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getCommit", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("getCommit", "projectId", projectId, "commitId", id)

	commit, _, err := conn.Commits.GetCommit(projectId, id, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			plugin.Logger(ctx).Warn("getCommit", "projectId", projectId, "commitId", id, "no project was found, returning empty result set")
			return nil, nil
		}
		plugin.Logger(ctx).Error("getCommit", "projectId", projectId, "commitId", id, "error", err)
		return nil, fmt.Errorf("unable to obtain commits for project_id %d\n%v", projectId, err)
	}

	commit.ProjectID = projectId
	commit.Message = strings.TrimRight(commit.Message, "\n") // remove trailing newline from commit message.

	plugin.Logger(ctx).Debug("getCommit", "completed successfully")
	return commit, nil
}

// Column Function
func commitColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "The ID (commit hash) of the commit.",
		},
		{
			Name:        "short_id",
			Type:        proto.ColumnType_STRING,
			Description: "The short ID (short commit hash) of the commit.",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the commit.",
		},
		{
			Name:        "author_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the commit author.",
		},
		{
			Name:        "author_email",
			Type:        proto.ColumnType_STRING,
			Description: "The email of the commit author.",
		},
		{
			Name:        "authored_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of commit.",
		},
		{
			Name:        "committer_name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the committer.",
		},
		{
			Name:        "committer_email",
			Type:        proto.ColumnType_STRING,
			Description: "The email address of the committer.",
		},
		{
			Name:        "committed_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of the commit.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp of the creation of commit.",
		},
		{
			Name:        "message",
			Type:        proto.ColumnType_STRING,
			Description: "The commit message.",
		},
		{
			Name:        "parent_ids",
			Type:        proto.ColumnType_JSON,
			Description: "Array of parent commit hashes.",
			Transform:   transform.FromField("ParentIDs"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project containing the commit - link to `gitlab_project.ID`",
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the commit.",
			Transform:   transform.FromField("WebURL"),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_STRING,
			Description: "Build state of the commit",
		},
		// Commit Stats
		{
			Name:        "commit_additions",
			Type:        proto.ColumnType_INT,
			Description: "Number of additions made in the commit",
			Transform:   transform.FromField("Stats.Additions"),
		},
		{
			Name:        "commit_deletions",
			Type:        proto.ColumnType_INT,
			Description: "Number of deletions made in the commit",
			Transform:   transform.FromField("Stats.Deletions"),
		},
		{
			Name:        "commit_total_changes",
			Type:        proto.ColumnType_INT,
			Description: "Total number of changes made in the commit",
			Transform:   transform.FromField("Stats.Total"),
		},
		// Pipeline Info
		{
			Name:        "pipeline_id",
			Type:        proto.ColumnType_INT,
			Description: "Identifier for the last pipeline instance triggered against the commit",
			Transform:   transform.FromField("LastPipeline.ID"),
		},
		{
			Name:        "pipeline_status",
			Type:        proto.ColumnType_STRING,
			Description: "Status of the last pipeline instance triggered against the commit",
			Transform:   transform.FromField("LastPipeline.Status"),
		},
		{
			Name:        "pipeline_source",
			Type:        proto.ColumnType_STRING,
			Description: "Source associated with the pipeline instance",
			Transform:   transform.FromField("LastPipeline.Source"),
		},
		{
			Name:        "pipeline_ref",
			Type:        proto.ColumnType_STRING,
			Description: "The ref that the pipeline was run against",
			Transform:   transform.FromField("LastPipeline.Ref"),
		},
		{
			Name:        "pipeline_sha",
			Type:        proto.ColumnType_STRING,
			Description: "The SHA of the commit the last pipeline instance was run against",
			Transform:   transform.FromField("LastPipeline.SHA"),
		},
		{
			Name:        "pipeline_url",
			Type:        proto.ColumnType_STRING,
			Description: "The URL of the pipeline in the web interface",
			Transform:   transform.FromField("LastPipeline.WebURL"),
		},
		{
			Name:        "pipeline_created",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp indicating when the last pipeline instance was created.",
			Transform:   transform.FromField("LastPipeline.CreatedAt"),
		},
		{
			Name:        "pipeline_updated",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp indicating when the last pipeline instance was updated.",
			Transform:   transform.FromField("LastPipeline.UpdatedAt"),
		},
	}
}
