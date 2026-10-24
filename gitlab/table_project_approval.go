package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableProjectApproval() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_approval",
		Description: "Obtain the merge request approval configuration for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectApproval,
		},
		Columns: projectApprovalColumns(),
	}
}

// Hydrate Functions
func listProjectApproval(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectApproval", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectApproval", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	plugin.Logger(ctx).Debug("listProjectApproval", "projectId", projectId)

	approval, _, err := conn.Projects.GetApprovalConfiguration(projectId)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectApproval", "projectId", projectId, "error", err)
		return nil, fmt.Errorf("unable to obtain approval configuration for project_id %d\n%v", projectId, err)
	}

	if approval == nil {
		plugin.Logger(ctx).Debug("listProjectApproval", "no approval configuration", projectId)
		return nil, nil
	}

	d.StreamListItem(ctx, approval)

	plugin.Logger(ctx).Debug("listProjectApproval", "completed successfully")
	return nil, nil
}

// Column Function
func projectApprovalColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "approvals_before_merge",
			Type:        proto.ColumnType_INT,
			Description: "The number of approvals required before a merge request can be merged.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "reset_approvals_on_push",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if existing approvals are reset when new commits are pushed.",
			Transform:   transform.FromField("ResetApprovalsOnPush"),
		},
		{
			Name:        "disable_overriding_approvers_per_merge_request",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if overriding the list of approvers per merge request is disabled.",
			Transform:   transform.FromField("DisableOverridingApproversPerMergeRequest"),
		},
		{
			Name:        "merge_requests_author_approval",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the author of a merge request is allowed to approve it.",
			Transform:   transform.FromField("MergeRequestsAuthorApproval"),
		},
		{
			Name:        "merge_requests_disable_committers_approval",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if committers to a merge request are prevented from approving it.",
			Transform:   transform.FromField("MergeRequestsDisableCommittersApproval"),
		},
		{
			Name:        "require_password_to_approve",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if approvers must re-authenticate with their password before approving.",
			Transform:   transform.FromField("RequirePasswordToApprove"),
		},
		{
			Name:        "selective_code_owner_removals",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if code owner approvals are reset only for the files changed.",
			Transform:   transform.FromField("SelectiveCodeOwnerRemovals"),
		},
		{
			Name:        "approvers",
			Type:        proto.ColumnType_JSON,
			Description: "The users configured as approvers at the project level.",
		},
		{
			Name:        "approver_groups",
			Type:        proto.ColumnType_JSON,
			Description: "The groups configured as approvers at the project level.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
