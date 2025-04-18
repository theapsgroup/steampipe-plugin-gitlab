package gitlab

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	api "gitlab.com/gitlab-org/api/client-go"
)

func tableGroupVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_group_variable",
		Description: "Obtain information on group level variables for a specific group in the GitLab instance.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "group_id",
					Require: plugin.Required,
				},
			},
			Hydrate: listGroupVars,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"group_id", "key"}),
			Hydrate:    getGroupVar,
		},
		Columns: groupVarColumns(),
	}
}

// Hydrate Functions
func listGroupVars(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listGroupVars", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listGroupVars", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	opt := &api.ListGroupVariablesOptions{
		Page:    1,
		PerPage: 50,
	}

	for {
		plugin.Logger(ctx).Debug("listGroupVars", "groupId", groupId, "page", opt.Page, "perPage", opt.PerPage)

		vars, resp, err := conn.GroupVariables.ListVariables(groupId, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listGroupVars", "groupId", groupId, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain group level variables for group_id %d\n%v", groupId, err)
		}

		for _, v := range vars {
			d.StreamListItem(ctx, v)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				plugin.Logger(ctx).Debug("listGroupVars", "completed successfully")
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listGroupVars", "completed successfully")
	return nil, nil
}

func getGroupVar(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getGroupVar", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getGroupVar", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	groupId := int(d.EqualsQuals["group_id"].GetInt64Value())
	key := d.EqualsQuals["key"].GetStringValue()
	plugin.Logger(ctx).Debug("getGroupVar", "groupId", groupId, "key", key)

	options := &api.GetGroupVariableOptions{}
	v, _, err := conn.GroupVariables.GetVariable(groupId, key, options, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getGroupVar", "groupId", groupId, "key", key, "error", err)
		return nil, fmt.Errorf("unable to obtain group level variable %s for group_id %d\n%v", key, groupId, err)
	}

	plugin.Logger(ctx).Debug("getGroupVar", "completed successfully")
	return v, nil
}

// Column Function
func groupVarColumns() []*plugin.Column {
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
			Name:        "group_id",
			Type:        proto.ColumnType_INT,
			Description: "The ID of the group this repository belongs to - link `gitlab_group.id`.",
			Transform:   transform.FromQual("group_id"),
		},
	}
}
