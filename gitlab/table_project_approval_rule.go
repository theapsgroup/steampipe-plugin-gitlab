package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectApprovalRule() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_approval_rule",
		Description: "Obtain information about merge request approval rules for a specific project within the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listProjectApprovalRules,
		},
		Columns: projectApprovalRuleColumns(),
	}
}

// Hydrate Functions
func listProjectApprovalRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listProjectApprovalRules", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectApprovalRules", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.GetProjectApprovalRulesListsOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectApprovalRules", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		rules, resp, err := conn.Projects.GetProjectApprovalRules(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectApprovalRules", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain approval rules for project_id %d\n%v", projectId, err)
		}

		for _, rule := range rules {
			d.StreamListItem(ctx, rule)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectApprovalRules", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectApprovalRules", "completed successfully")
	return nil, nil
}

// Column Function
func projectApprovalRuleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the approval rule.",
		},
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "The name of the approval rule.",
		},
		{
			Name:        "rule_type",
			Type:        proto.ColumnType_STRING,
			Description: "The type of the approval rule (e.g. regular, code_owner, any_approver).",
		},
		{
			Name:        "approvals_required",
			Type:        proto.ColumnType_INT,
			Description: "The number of approvals required for this rule.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "eligible_approvers",
			Type:        proto.ColumnType_JSON,
			Description: "The users eligible to approve under this rule (resolved from users and groups).",
		},
		{
			Name:        "users",
			Type:        proto.ColumnType_JSON,
			Description: "The users explicitly added as approvers on this rule.",
		},
		{
			Name:        "groups",
			Type:        proto.ColumnType_JSON,
			Description: "The groups explicitly added as approvers on this rule.",
		},
		{
			Name:        "contains_hidden_groups",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the rule contains groups the requesting user has no access to.",
			Transform:   transform.FromField("ContainsHiddenGroups"),
		},
		{
			Name:        "protected_branches",
			Type:        proto.ColumnType_JSON,
			Description: "The protected branches this rule applies to.",
		},
		{
			Name:        "applies_to_all_protected_branches",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the rule applies to all protected branches.",
			Transform:   transform.FromField("AppliesToAllProtectedBranches"),
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The project id - link to `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
