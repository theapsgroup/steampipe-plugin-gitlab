package gitlab

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "github.com/xanzy/go-gitlab"
)

func tableProjectVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_variable",
		Description: "Variables for a GitLab Project",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "project_id",
					Require: plugin.Required,
				},
			},
			Hydrate: listProjectVars,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "key"}),
			Hydrate:    getProjectVar,
		},
		Columns: projectVarColumns(),
	}
}

func listProjectVars(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())

	opt := &api.ListProjectVariablesOptions{
		Page:    1,
		PerPage: 20,
	}
	for {
		vars, resp, err := conn.ProjectVariables.ListVariables(projectId, opt)
		if err != nil {
			return nil, err
		}

		for _, v := range vars {
			d.StreamListItem(ctx, v)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func getProjectVar(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	key := d.EqualsQuals["key"].GetStringValue()
	opt := &api.GetProjectVariableOptions{}

	v, _, err := conn.ProjectVariables.GetVariable(projectId, key, opt)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func projectVarColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "key",
			Type:        proto.ColumnType_STRING,
			Description: "The key of the variable.",
		},
		{
			Name:        "value",
			Type:        proto.ColumnType_STRING,
			Description: "The value of the variable.",
		},
		{
			Name:        "variable_type",
			Type:        proto.ColumnType_STRING,
			Description: "The type of the variable (env var, etc).",
		},
		{
			Name:        "environment_scope",
			Type:        proto.ColumnType_STRING,
			Description: "The environment(s) that this variable is in scope for.",
		},
		{
			Name:        "protected",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the variable is only applied to protected branches.",
		},
		{
			Name:        "masked",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the variable is masked (hidden) in job logs.",
		},
		{
			Name:        "raw",
			Type:        proto.ColumnType_BOOL,
			Description: "Indicates if the variable is is a raw format.",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the project this repository belongs to - link `gitlab_project.id`.",
			Transform:   transform.FromQual("project_id"),
		},
	}
}
