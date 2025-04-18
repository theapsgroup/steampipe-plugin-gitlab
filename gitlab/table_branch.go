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

func tableBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_branch",
		Description: "Obtain information on branches for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listBranches,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "name"}),
			Hydrate:    getBranch,
		},
		Columns: branchColumns(),
	}
}

// Hydrate Functions
func listBranches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listBranches", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listBranches", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListBranchesOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	for {
		plugin.Logger(ctx).Debug("listBranches", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)

		branches, resp, err := conn.Branches.ListBranches(projectId, opt)
		if err != nil {
			// Handle error of project id not being valid.
			if strings.Contains(err.Error(), "404") {
				plugin.Logger(ctx).Warn("listBranches", "projectId", projectId, "no project was found, returning empty result set")
				return nil, nil
			}
			plugin.Logger(ctx).Error("listBranches", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain branches for project_id %d\n%v", projectId, err)
		}

		for _, branch := range branches {
			d.StreamListItem(ctx, branch)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listBranches", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listBranches", "completed successfully")
	return nil, nil
}

func getBranch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getBranch", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getBranch", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	name := d.EqualsQuals["name"].GetStringValue()
	plugin.Logger(ctx).Debug("getBranch", "projectId", projectId, "name", name)

	branch, _, err := conn.Branches.GetBranch(projectId, name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			plugin.Logger(ctx).Warn("getBranch", "projectId", projectId, "name", name, "no project was found, returning empty result set")
			return nil, nil
		}
		plugin.Logger(ctx).Error("getBranch", "projectId", projectId, "name", name, "error", err)
		return nil, fmt.Errorf("unable to obtain branch %s for project_id %d\n%v", name, projectId, err)
	}

	plugin.Logger(ctx).Debug("getBranch", "completed successfully")
	return branch, nil
}

// Column Function
func branchColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project containing the branches - link to `gitlab_project.ID`",
			Transform:   transform.FromQual("project_id"),
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the branch.",
		},
		{
			Name:        "protected",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the branch is protected or not.",
		},
		{
			Name:        "merged",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the branch has been merged into the trunk.",
		},
		{
			Name:        "default",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the branch is the default branch of the project.",
		},
		{
			Name:        "can_push",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the current user can push to this branch.",
		},
		{
			Name:        "devs_can_push",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users with the `developer` level of access can push to the branch.",
			Transform:   transform.FromField("DevelopersCanPush"),
		},
		{
			Name:        "devs_can_merge",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if users with the `developer` level of access can merge the branch.",
			Transform:   transform.FromField("DevelopersCanMerge"),
		},
		{
			Name:        "web_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the branch.",
			Transform:   transform.FromField("WebURL").NullIfZero(),
		},
		{
			Name:        "commit_id",
			Type:        proto.ColumnType_STRING,
			Description: "The latest commit hash on the branch.",
			Transform:   transform.FromField("Commit.ID"),
		},
		{
			Name:        "commit_short_id",
			Type:        proto.ColumnType_STRING,
			Description: "The latest short commit hash on the branch.",
			Transform:   transform.FromField("Commit.ShortID"),
		},
		{
			Name:        "commit_title",
			Type:        proto.ColumnType_STRING,
			Description: "The title of the latest commit on the branch.",
			Transform:   transform.FromField("Commit.Title"),
		},
		{
			Name:        "commit_message",
			Type:        proto.ColumnType_STRING,
			Description: "The commit message associated with the latest commit on the branch.",
			Transform:   transform.FromField("Commit.Message"),
		},
		{
			Name:        "commit_email",
			Type:        proto.ColumnType_STRING,
			Description: "The email address associated with the latest commit on the branch.",
			Transform:   transform.FromField("Commit.CommitterEmail"),
		},
		{
			Name:        "commit_date",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The date of the latest commit on the branch.",
			Transform:   transform.FromField("Commit.CommittedDate"),
		},
		{
			Name:        "commit_url",
			Type:        proto.ColumnType_STRING,
			Description: "The url of the commit on the branch.",
			Transform:   transform.FromField("Commit.WebURL"),
		},
	}
}
