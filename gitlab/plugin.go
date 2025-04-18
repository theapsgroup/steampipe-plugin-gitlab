package gitlab

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-gitlab",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: ignoreErrorsContaining([]string{
				"404",
			})},
		TableMap: map[string]*plugin.Table{
			"gitlab_application":                tableApplication(),
			"gitlab_branch":                     tableBranch(),
			"gitlab_commit":                     tableCommit(),
			"gitlab_epic":                       tableEpic(),
			"gitlab_group":                      tableGroup(),
			"gitlab_group_access_request":       tableGroupAccessRequest(),
			"gitlab_group_hook":                 tableGroupHook(),
			"gitlab_group_iteration":            tableGroupIteration(),
			"gitlab_group_member":               tableGroupMember(),
			"gitlab_group_project":              tableGroupProject(),
			"gitlab_group_push_rule":            tableGroupPushRule(),
			"gitlab_group_subgroup":             tableGroupSubgroup(),
			"gitlab_group_variable":             tableGroupVariable(),
			"gitlab_instance_variable":          tableInstanceVariable(),
			"gitlab_issue":                      tableIssue(),
			"gitlab_merge_request":              tableMergeRequest(),
			"gitlab_merge_request_diffs":        tableListMergeRequestDiffs(),
			"gitlab_my_event":                   tableMyEvents(),
			"gitlab_my_issue":                   tableMyIssue(),
			"gitlab_my_project":                 tableMyProject(),
			"gitlab_project":                    tableProject(),
			"gitlab_project_access_request":     tableProjectAccessRequest(),
			"gitlab_project_container_registry": tableProjectContainerRegistry(),
			"gitlab_project_deployment":         tableProjectDeployment(),
			"gitlab_project_iteration":          tableProjectIteration(),
			"gitlab_project_job":                tableProjectJob(),
			"gitlab_project_member":             tableProjectMember(),
			"gitlab_project_pages_domain":       tableProjectPagesDomain(),
			"gitlab_project_pipeline":           tableProjectPipeline(),
			"gitlab_project_pipeline_detail":    tableProjectPipelineDetail(),
			"gitlab_project_protected_branch":   tableProjectProtectedBranch(),
			"gitlab_project_repository":         tableProjectRepository(),
			"gitlab_project_repository_file":    tableProjectRepositoryFile(),
			"gitlab_project_variable":           tableProjectVariable(),
			"gitlab_setting":                    tableSetting(),
			"gitlab_snippet":                    tableSnippet(),
			"gitlab_user":                       tableUser(),
			"gitlab_user_event":                 tableUserEvents(),
			"gitlab_version":                    tableVersion(),
		},
	}

	return p
}

func ignoreErrorsContaining(terms []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err != nil {
			for _, term := range terms {
				if strings.Contains(err.Error(), term) {
					return true
				}
			}
		}
		return false
	}
}
