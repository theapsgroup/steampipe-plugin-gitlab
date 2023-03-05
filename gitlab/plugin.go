package gitlab

import (
	"context"

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
		TableMap: map[string]*plugin.Table{
			"gitlab_version":                    tableVersion(),
			"gitlab_user":                       tableUser(),
			"gitlab_group":                      tableGroup(),
			"gitlab_project":                    tableProject(),
			"gitlab_issue":                      tableIssue(),
			"gitlab_branch":                     tableBranch(),
			"gitlab_commit":                     tableCommit(),
			"gitlab_merge_request":              tableMergeRequest(),
			"gitlab_group_member":               tableGroupMember(),
			"gitlab_project_member":             tableProjectMember(),
			"gitlab_snippet":                    tableSnippet(),
			"gitlab_project_pipeline":           tableProjectPipeline(),
			"gitlab_project_pipeline_detail":    tableProjectPipelineDetail(),
			"gitlab_project_job":                tableProjectJob(),
			"gitlab_project_repository":         tableProjectRepository(),
			"gitlab_project_repository_file":    tableProjectRepositoryFile(),
			"gitlab_my_project":                 tableMyProject(),
			"gitlab_my_issue":                   tableMyIssue(),
			"gitlab_epic":                       tableEpic(),
			"gitlab_group_iteration":            tableGroupIteration(),
			"gitlab_project_iteration":          tableProjectIteration(),
			"gitlab_group_push_rule":            tableGroupPushRule(),
			"gitlab_group_hook":                 tableGroupHook(),
			"gitlab_project_protected_branch":   tableProjectProtectedBranch(),
			"gitlab_project_pages_domain":       tableProjectPagesDomain(),
			"gitlab_setting":                    tableSetting(),
			"gitlab_merge_request_change":       tableMergeRequestChange(),
			"gitlab_project_variable":           tableProjectVariable(),
			"gitlab_group_variable":             tableGroupVariable(),
			"gitlab_instance_variable":          tableInstanceVariable(),
			"gitlab_group_access_request":       tableGroupAccessRequest(),
			"gitlab_project_access_request":     tableProjectAccessRequest(),
			"gitlab_application":                tableApplication(),
			"gitlab_project_container_registry": tableProjectContainerRegistry(),
			"gitlab_project_deployment":         tableProjectDeployment(),
		},
	}

	return p
}
