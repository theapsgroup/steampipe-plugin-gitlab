package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectProtectedBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_protected_branch",
		Description: "Protected Branches for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectProtectedBranches,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_INT,
				Description: "The ID of the protected branch.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the protected branch.",
			},
			{
				Name:        "allow_force_push",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if force pushing is allowed on the protected branch.",
			},
			{
				Name:        "code_owner_approval_required",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates if code owner approval is required.",
			},
			{
				Name:        "push_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Array of push access levels.",
			},
			{
				Name:        "merge_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Array of merge access levels.",
			},
			{
				Name:        "unprotect_access_levels",
				Type:        proto.ColumnType_JSON,
				Description: "Array of unprotected access levels.",
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_INT,
				Description: "The ID of the project the protected branch belongs to - link `gitlab_project.id`.",
				Transform:   transform.FromQual("project_id"),
			},
		},
	}
}

func listProjectProtectedBranches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.KeyColumnQuals["project_id"].GetInt64Value())
	opt := &api.ListProtectedBranchesOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		branches, resp, err := conn.ProtectedBranches.ListProtectedBranches(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, branch := range branches {
			d.StreamListItem(ctx, branch)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
