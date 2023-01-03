package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
	"strings"
	"time"
)

type Branch struct {
	ProjectID          int
	Name               string
	Protected          bool
	Merged             bool
	Default            bool
	CanPush            bool
	DevelopersCanPush  bool
	DevelopersCanMerge bool
	WebUrl             string
	CommitID           string
	CommitShortID      string
	CommitTitle        string
	CommitEmail        string
	CommitDate         *time.Time
}

func tableBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_branch",
		Description: "Branches in the given project.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listBranches,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "name"}),
			Hydrate:    getBranch,
		},
		Columns: []*plugin.Column{
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project containing the branches - link to `gitlab_project.ID`"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the branch."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Description: "Indicates if the branch is protected or not."},
			{Name: "merged", Type: proto.ColumnType_BOOL, Description: "Indicates if the branch has been merged into the trunk."},
			{Name: "default", Type: proto.ColumnType_BOOL, Description: "Indicates if the branch is the default branch of the project."},
			{Name: "can_push", Type: proto.ColumnType_BOOL, Description: "Indicates if the current user can push to this branch."},
			{Name: "devs_can_push", Type: proto.ColumnType_BOOL, Description: "Indicates if users with the `developer` level of access can push to the branch.", Transform: transform.FromField("DevelopersCanPush")},
			{Name: "devs_can_merge", Type: proto.ColumnType_BOOL, Description: "Indicates if users with the `developer` level of access can merge the branch.", Transform: transform.FromField("DevelopersCanMerge")},
			{Name: "web_url", Type: proto.ColumnType_STRING, Description: "The url of the branch.", Transform: transform.FromField("WebUrl").NullIfZero()},
			{Name: "commit_id", Type: proto.ColumnType_STRING, Description: "The latest commit hash on the branch."},
			{Name: "commit_short_id", Type: proto.ColumnType_STRING, Description: "The latest short commit hash on the branch."},
			{Name: "commit_title", Type: proto.ColumnType_STRING, Description: "The title of the latest commit on the branch."},
			{Name: "commit_email", Type: proto.ColumnType_STRING, Description: "The email address associated with the latest commit on the branch."},
			{Name: "commit_date", Type: proto.ColumnType_TIMESTAMP, Description: "The date of the latest commit on the branch."},
		},
	}
}

func listBranches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opt := &api.ListBranchesOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 10,
	}}

	for {
		branches, resp, err := conn.Branches.ListBranches(projectId, opt)
		if err != nil {
			// Handle error of project id not being valid.
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			return nil, err
		}

		for _, branch := range branches {
			d.StreamListItem(ctx, &Branch{
				ProjectID:          projectId,
				Name:               branch.Name,
				Protected:          branch.Protected,
				Merged:             branch.Merged,
				Default:            branch.Default,
				CanPush:            branch.CanPush,
				DevelopersCanPush:  branch.DevelopersCanPush,
				DevelopersCanMerge: branch.DevelopersCanMerge,
				WebUrl:             branch.WebURL,
				CommitID:           branch.Commit.ID,
				CommitShortID:      branch.Commit.ShortID,
				CommitTitle:        branch.Commit.Title,
				CommitEmail:        branch.Commit.CommitterEmail,
				CommitDate:         branch.Commit.CommittedDate,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func getBranch(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	name := d.EqualsQuals["name"].GetStringValue()

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	branch, _, err := conn.Branches.GetBranch(projectId, name)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}

	return &Branch{
		ProjectID:          projectId,
		Name:               branch.Name,
		Protected:          branch.Protected,
		Merged:             branch.Merged,
		Default:            branch.Default,
		CanPush:            branch.CanPush,
		DevelopersCanPush:  branch.DevelopersCanPush,
		DevelopersCanMerge: branch.DevelopersCanMerge,
		WebUrl:             branch.WebURL,
		CommitID:           branch.Commit.ID,
		CommitShortID:      branch.Commit.ShortID,
		CommitTitle:        branch.Commit.Title,
		CommitEmail:        branch.Commit.CommitterEmail,
		CommitDate:         branch.Commit.CommittedDate,
	}, nil
}
