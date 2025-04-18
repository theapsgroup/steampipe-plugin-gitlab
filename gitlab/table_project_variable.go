package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableProjectVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_project_variable",
		Description: "Obtain information on project level variables for a specific group in the GitLab instance.",
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
	plugin.Logger(ctx).Debug("listProjectVars", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listProjectVars", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	opt := &api.ListProjectVariablesOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listProjectVars", "projectId", projectId, "page", opt.Page, "perPage", opt.PerPage)
		vars, resp, err := conn.ProjectVariables.ListVariables(projectId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listProjectVars", "projectId", projectId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain issues for project_id %d\n%v", projectId, err)
		}

		for _, v := range vars {
			d.StreamListItem(ctx, v)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listProjectVars", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listProjectVars", "completed successfully")
	return nil, nil
}

func getProjectVar(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getProjectVar", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getProjectVar", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	projectId := int(d.EqualsQuals["project_id"].GetInt64Value())
	key := d.EqualsQuals["key"].GetStringValue()
	opt := &api.GetProjectVariableOptions{}

	plugin.Logger(ctx).Debug("getProjectVar", "projectId", projectId, "key", key)
	v, _, err := conn.ProjectVariables.GetVariable(projectId, key, opt)
	if err != nil {
		plugin.Logger(ctx).Error("getProjectVar", "projectId", projectId, "key", key, "error", err)
		return nil, fmt.Errorf("unable to obtain variable %s for project_id %d\n%v", key, projectId, err)
	}

	plugin.Logger(ctx).Debug("getProjectVar", "completed successfully")
	return v, nil
}

// Column Function
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
