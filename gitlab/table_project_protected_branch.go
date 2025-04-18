package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectProtectedBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_protected_branch",
		Description: "Obtain information about protected branches for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectProtectedBranches,
		},
		Columns: protectedBranchColumns(),
	}
}

// Hydration Functions
func listProjectProtectedBranches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectProtectedBranches", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectProtectedBranches", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProtectedBranchesOptions{
		ListOptions: api.ListOptions{
			Page:    1,
			PerPage: 50,
		},
	}

	for {
		plugin.Logger(ctx).Debug("listProjectProtectedBranches", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		branches, resp, err := conn.ProtectedBranches.ListProtectedBranches(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectProtectedBranches", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain protected branches for project_id %d\n%v", projectId, err)
		}

		for _, branch := range branches {
			d.StreamListItem(ctx, branch)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectProtectedBranches", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectProtectedBranches", "completed successfully")
	return nil, nil
}

// Column Function
func protectedBranchColumns() []*plugin.Column {
	return []*plugin.Column{
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
	}
}
