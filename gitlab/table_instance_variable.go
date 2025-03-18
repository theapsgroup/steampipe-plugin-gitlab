package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableInstanceVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_instance_variable",
		Description: "Obtain information on instance level variables within the GitLab instance (only available in self-hosted model).",
		List: &plugin.ListConfig{
			Hydrate: listInstanceVars,
		},
		Columns: instanceVarColumns(),
	}
}

// Hydrate Functions
func listInstanceVars(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listInstanceVars", "started")
	// Not available on public, only self-hosted.
	if isPublicGitLab(d) {
		plugin.Logger(ctx).Warn("listInstanceVars", "non-self hosted instance - exiting with empty result set")
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listInstanceVars", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	opt := &api.ListInstanceVariablesOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listInstanceVars", "page", opt.Page, "perPage", opt.PerPage)

		vars, resp, err := conn.InstanceVariables.ListVariables(opt)
		if err != nil {
			plugin.Logger(ctx).Error("listInstanceVars", "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain instance level variables\n%v", err)
		}

		for _, v := range vars {
			d.StreamListItem(ctx, v)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listInstanceVars", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listInstanceVars", "completed successfully")
	return nil, nil
}

// Column Function
func instanceVarColumns() []*plugin.Column {
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
	}
}
