package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			"gitlab_version":                 tableVersion(),
			"gitlab_user":                    tableUser(),
			"gitlab_group":                   tableGroup(),
			"gitlab_project":                 tableProject(),
			"gitlab_issue":                   tableIssue(),
			"gitlab_branch":                  tableBranch(),
			"gitlab_commit":                  tableCommit(),
			"gitlab_merge_request":           tableMergeRequest(),
			"gitlab_project_merge_request":   tableProjectMergeRequest(),
			"gitlab_group_member":            tableGroupMember(),
			"gitlab_project_member":          tableProjectMember(),
			"gitlab_snippet":                 tableSnippet(),
			"gitlab_project_pipeline":        tableProjectPipeline(),
			"gitlab_project_pipeline_detail": tableProjectPipelineDetail(),
		},
	}

	return p
}